//-----------------------------------------------------------------------------
/*

Complex 3D CAD Architecture Example - Parametric Robotic Joint Assembly

This example demonstrates best practices for organizing complex 3D CAD models:

Architecture Overview:
1. Configuration Management - All parameters in structs
2. Component Library - Reusable mechanical parts
3. Sub-assemblies - Logical grouping of related components
4. Assembly Logic - Main assembly with proper transforms
5. Variant Generation - Multiple configurations from same code

Components:
- Base mount plate with mounting holes
- Bearing housings (top and bottom)
- Drive shaft with keyway
- Input and output gears
- Cover plate
- Fasteners (bolts, washers)

*/
//-----------------------------------------------------------------------------

package main

import (
    "log"
    "math"

    "github.com/deadsy/sdfx/obj"
    "github.com/deadsy/sdfx/render"
    "github.com/deadsy/sdfx/sdf"
    v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// Configuration Management
//-----------------------------------------------------------------------------

// MaterialConfig defines material properties and shrinkage factors
type MaterialConfig struct {
    Name          string
    ShrinkageFactor float64
    MinWallThickness float64
    GeneralRounding  float64
}

// StandardMaterials provides common 3D printing materials
var StandardMaterials = map[string]MaterialConfig{
    "PLA": {
        Name:             "PLA",
        ShrinkageFactor:  1.0 / 0.998,
        MinWallThickness: 1.2,
        GeneralRounding:  0.5,
    },
    "ABS": {
        Name:             "ABS",
        ShrinkageFactor:  1.0 / 0.995,
        MinWallThickness: 1.5,
        GeneralRounding:  0.5,
    },
    "PETG": {
        Name:             "PETG",
        ShrinkageFactor:  1.0 / 0.997,
        MinWallThickness: 1.2,
        GeneralRounding:  0.4,
    },
}

// JointConfig defines all parameters for the robotic joint
type JointConfig struct {
    // Material settings
    Material MaterialConfig
    
    // Base dimensions
    BaseDiameter     float64
    BaseThickness    float64
    BaseHoleCount    int
    BaseHoleDiameter float64
    BaseMountRadius  float64
    
    // Shaft parameters
    ShaftDiameter    float64
    ShaftLength      float64
    KeywayWidth      float64
    KeywayDepth      float64
    
    // Bearing parameters
    BearingOD        float64
    BearingID        float64
    BearingThickness float64
    BearingClearance float64
    
    // Housing parameters
    HousingWallThickness float64
    HousingFlangeHeight  float64
    
    // Gear parameters
    GearModule       float64
    InputTeeth       int
    OutputTeeth      int
    GearThickness    float64
    GearPressureAngle float64
    
    // Fastener parameters
    BoltDiameter     float64
    BoltHeadDiameter float64
    BoltHeadHeight   float64
    
    // Assembly clearances
    BoltClearance    float64
    GeneralClearance float64
    
    // Quality settings
    MeshResolution   int
}

// NewDefaultJointConfig creates a standard configuration
func NewDefaultJointConfig() *JointConfig {
    return &JointConfig{
        Material:             StandardMaterials["PLA"],
        BaseDiameter:         80.0,
        BaseThickness:        8.0,
        BaseHoleCount:        6,
        BaseHoleDiameter:     5.0,
        BaseMountRadius:      30.0,
        ShaftDiameter:        12.0,
        ShaftLength:          50.0,
        KeywayWidth:          3.0,
        KeywayDepth:          1.5,
        BearingOD:            32.0,
        BearingID:            12.0,
        BearingThickness:     10.0,
        BearingClearance:     0.2,
        HousingWallThickness: 4.0,
        HousingFlangeHeight:  3.0,
        GearModule:           1.5,
        InputTeeth:           20,
        OutputTeeth:          40,
        GearThickness:        8.0,
        GearPressureAngle:    sdf.DtoR(20.0),
        BoltDiameter:         3.0,
        BoltHeadDiameter:     5.5,
        BoltHeadHeight:       2.0,
        BoltClearance:        0.3,
        GeneralClearance:     0.2,
        MeshResolution:       300,
    }
}

//-----------------------------------------------------------------------------
// Component Library - Reusable Parts
//-----------------------------------------------------------------------------

// ComponentMode defines what type of feature to generate
type ComponentMode int

const (
    ModeBody ComponentMode = iota
    ModeHole
    ModeBoss
    ModeClearance
)

// BaseMount creates the mounting plate for the joint
func BaseMount(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    switch mode {
    case ModeBody:
        // Create main circular base
        base, err := sdf.Cylinder3D(cfg.BaseThickness, cfg.BaseDiameter/2, cfg.Material.GeneralRounding)
        if err != nil {
            return nil, err
        }
        return base, nil
        
    case ModeHole:
        // Center shaft hole
        centerHole, err := sdf.Cylinder3D(cfg.BaseThickness*1.5, cfg.ShaftDiameter/2+cfg.GeneralClearance, 0)
        if err != nil {
            return nil, err
        }
        
        // Mounting bolt holes
        boltHole, err := sdf.Cylinder3D(cfg.BaseThickness*1.5, cfg.BaseHoleDiameter/2, 0)
        if err != nil {
            return nil, err
        }
        
        // Position bolt holes in circular pattern
        var boltHoles []sdf.SDF3
        angleStep := 2 * math.Pi / float64(cfg.BaseHoleCount)
        for i := 0; i < cfg.BaseHoleCount; i++ {
            angle := float64(i) * angleStep
            x := cfg.BaseMountRadius * math.Cos(angle)
            y := cfg.BaseMountRadius * math.Sin(angle)
            hole := sdf.Transform3D(boltHole, sdf.Translate3d(v3.Vec{x, y, 0}))
            boltHoles = append(boltHoles, hole)
        }
        
        // Combine all holes
        allHoles := append([]sdf.SDF3{centerHole}, boltHoles...)
        return sdf.Union3D(allHoles...), nil
        
    default:
        return nil, nil
    }
}

// BearingHousing creates the housing for a bearing
func BearingHousing(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    housingOD := cfg.BearingOD + 2*cfg.HousingWallThickness
    housingHeight := cfg.BearingThickness + cfg.HousingFlangeHeight
    
    switch mode {
    case ModeBody:
        // Main cylindrical body
        body, err := sdf.Cylinder3D(housingHeight, housingOD/2, cfg.Material.GeneralRounding)
        if err != nil {
            return nil, err
        }
        
        // Flange for mounting
        flangeHeight := cfg.HousingFlangeHeight
        flangeRadius := housingOD/2 + cfg.HousingWallThickness
        flange, err := sdf.Cylinder3D(flangeHeight, flangeRadius, cfg.Material.GeneralRounding)
        if err != nil {
            return nil, err
        }
        
        // Position flange at bottom
        zOffset := -(housingHeight - flangeHeight) / 2
        flange = sdf.Transform3D(flange, sdf.Translate3d(v3.Vec{0, 0, zOffset}))
        
        // Union with smooth blending
        combined := sdf.Union3D(body, flange)
        combined.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(cfg.Material.GeneralRounding))
        
        return combined, nil
        
    case ModeHole:
        // Bearing pocket
        bearingPocket, err := sdf.Cylinder3D(cfg.BearingThickness, cfg.BearingOD/2+cfg.BearingClearance, 0)
        if err != nil {
            return nil, err
        }
        
        // Center through hole
        throughHole, err := sdf.Cylinder3D(housingHeight*1.5, cfg.BearingID/2+cfg.GeneralClearance, 0)
        if err != nil {
            return nil, err
        }
        
        // Position bearing pocket
        zOffset := cfg.HousingFlangeHeight / 2
        bearingPocket = sdf.Transform3D(bearingPocket, sdf.Translate3d(v3.Vec{0, 0, zOffset}))
        
        return sdf.Union3D(bearingPocket, throughHole), nil
        
    default:
        return nil, nil
    }
}

// DriveShaft creates shaft with optional keyway
func DriveShaft(cfg *JointConfig, includeKeyway bool) (sdf.SDF3, error) {
    // Main shaft cylinder
    shaft, err := sdf.Cylinder3D(cfg.ShaftLength, cfg.ShaftDiameter/2, 0)
    if err != nil {
        return nil, err
    }
    
    if includeKeyway {
        // Create keyway slot
        keyway, err := sdf.Box3D(v3.Vec{cfg.ShaftDiameter, cfg.KeywayWidth, cfg.ShaftLength}, 0)
        if err != nil {
            return nil, err
        }
        
        // Position keyway at edge of shaft
        xOffset := cfg.ShaftDiameter/2 - cfg.KeywayDepth/2
        keyway = sdf.Transform3D(keyway, sdf.Translate3d(v3.Vec{xOffset, 0, 0}))
        
        // Subtract keyway from shaft
        shaft = sdf.Difference3D(shaft, keyway)
    }
    
    return shaft, nil
}

// ParametricGear creates a gear with the given configuration
func ParametricGear(cfg *JointConfig, numTeeth int, withBore bool) (sdf.SDF3, error) {
    // Create involute gear profile
    gearParams := obj.InvoluteGearParms{
        NumberTeeth:   numTeeth,
        Module:        cfg.GearModule,
        PressureAngle: cfg.GearPressureAngle,
        RingWidth:     cfg.Material.MinWallThickness * 2,
        Facets:        8,
    }
    
    gear2d, err := obj.InvoluteGear(&gearParams)
    if err != nil {
        return nil, err
    }
    
    // Extrude to 3D
    gear3d := sdf.Extrude3D(gear2d, cfg.GearThickness)
    
    if withBore {
        // Add center bore with keyway
        bore, err := sdf.Cylinder3D(cfg.GearThickness*1.5, cfg.ShaftDiameter/2+cfg.GeneralClearance, 0)
        if err != nil {
            return nil, err
        }
        
        // Keyway slot in gear
        keySlot, err := sdf.Box3D(v3.Vec{cfg.ShaftDiameter, cfg.KeywayWidth+cfg.GeneralClearance*2, cfg.GearThickness*1.5}, 0)
        if err != nil {
            return nil, err
        }
        xOffset := cfg.ShaftDiameter/2 - cfg.KeywayDepth/2
        keySlot = sdf.Transform3D(keySlot, sdf.Translate3d(v3.Vec{xOffset, 0, 0}))
        
        boreWithKey := sdf.Union3D(bore, keySlot)
        gear3d = sdf.Difference3D(gear3d, boreWithKey)
    }
    
    return gear3d, nil
}

// CoverPlate creates a protective cover for the assembly
func CoverPlate(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    plateRadius := cfg.BaseDiameter / 2
    plateThickness := cfg.BaseThickness * 0.6
    
    switch mode {
    case ModeBody:
        plate, err := sdf.Cylinder3D(plateThickness, plateRadius, cfg.Material.GeneralRounding)
        if err != nil {
            return nil, err
        }
        return plate, nil
        
    case ModeHole:
        // Center inspection hole
        inspectionHole, err := sdf.Cylinder3D(plateThickness*1.5, plateRadius*0.4, 0)
        if err != nil {
            return nil, err
        }
        
        // Ventilation slots
        slot, err := sdf.Box3D(v3.Vec{plateRadius * 0.6, 3.0, plateThickness * 1.5}, cfg.Material.GeneralRounding)
        if err != nil {
            return nil, err
        }
        
        var slots []sdf.SDF3
        for i := 0; i < 4; i++ {
            angle := float64(i) * math.Pi / 2
            m := sdf.RotateZ(angle)
            m = sdf.Translate3d(v3.Vec{0, plateRadius * 0.6, 0}).Mul(m)
            s := sdf.Transform3D(slot, m)
            slots = append(slots, s)
        }
        
        allHoles := append([]sdf.SDF3{inspectionHole}, slots...)
        return sdf.Union3D(allHoles...), nil
        
    default:
        return nil, nil
    }
}

// CounterboreHole creates a bolt hole with counterbore for bolt head
func CounterboreHole(cfg *JointConfig, depth float64) (sdf.SDF3, error) {
    // Main bolt shaft hole
    boltHole, err := sdf.Cylinder3D(depth, (cfg.BoltDiameter+cfg.BoltClearance)/2, 0)
    if err != nil {
        return nil, err
    }
    
    // Counterbore for bolt head
    cbDepth := cfg.BoltHeadHeight + 0.5
    counterbore, err := sdf.Cylinder3D(cbDepth, (cfg.BoltHeadDiameter+cfg.BoltClearance)/2, 0)
    if err != nil {
        return nil, err
    }
    
    // Position counterbore at top
    zOffset := (depth - cbDepth) / 2
    counterbore = sdf.Transform3D(counterbore, sdf.Translate3d(v3.Vec{0, 0, zOffset}))
    
    return sdf.Union3D(boltHole, counterbore), nil
}

//-----------------------------------------------------------------------------
// Sub-Assembly Functions
//-----------------------------------------------------------------------------

// BasePlateAssembly creates the complete base plate with all holes
func BasePlateAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // Get base body
    body, err := BaseMount(cfg, ModeBody)
    if err != nil {
        return nil, err
    }
    
    // Get all holes
    holes, err := BaseMount(cfg, ModeHole)
    if err != nil {
        return nil, err
    }
    
    // Combine
    base := sdf.Difference3D(body, holes)
    
    return base, nil
}

// BearingAssembly creates bearing housing with integrated features
func BearingAssembly(cfg *JointConfig, upperHousing bool) (sdf.SDF3, error) {
    // Housing body
    body, err := BearingHousing(cfg, ModeBody)
    if err != nil {
        return nil, err
    }
    
    // Housing holes
    holes, err := BearingHousing(cfg, ModeHole)
    if err != nil {
        return nil, err
    }
    
    housing := sdf.Difference3D(body, holes)
    
    // Add mounting bolt holes for upper housing
    if upperHousing {
        housingOD := cfg.BearingOD + 2*cfg.HousingWallThickness
        housingHeight := cfg.BearingThickness + cfg.HousingFlangeHeight
        mountRadius := (housingOD/2 + cfg.HousingWallThickness) * 0.8
        
        boltHole, err := CounterboreHole(cfg, housingHeight)
        if err != nil {
            return nil, err
        }
        
        var bolts []sdf.SDF3
        for i := 0; i < 4; i++ {
            angle := float64(i)*math.Pi/2 + math.Pi/4
            x := mountRadius * math.Cos(angle)
            y := mountRadius * math.Sin(angle)
            b := sdf.Transform3D(boltHole, sdf.Translate3d(v3.Vec{x, y, 0}))
            bolts = append(bolts, b)
        }
        
        housing = sdf.Difference3D(housing, sdf.Union3D(bolts...))
    }
    
    return housing, nil
}

// DriveTrainAssembly creates the complete gear and shaft assembly
func DriveTrainAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // Create shaft
    shaft, err := DriveShaft(cfg, true)
    if err != nil {
        return nil, err
    }
    
    // Create input gear (smaller)
    inputGear, err := ParametricGear(cfg, cfg.InputTeeth, true)
    if err != nil {
        return nil, err
    }
    
    // Create output gear (larger)
    outputGear, err := ParametricGear(cfg, cfg.OutputTeeth, false)
    if err != nil {
        return nil, err
    }
    
    // Calculate gear positioning
    pitchRadius1 := float64(cfg.InputTeeth) * cfg.GearModule / 2
    pitchRadius2 := float64(cfg.OutputTeeth) * cfg.GearModule / 2
    centerDistance := pitchRadius1 + pitchRadius2
    
    // Position input gear on shaft
    inputGear = sdf.Transform3D(inputGear, sdf.Translate3d(v3.Vec{0, 0, cfg.ShaftLength/4}))
    
    // Position output gear offset
    outputGear = sdf.Transform3D(outputGear, sdf.Translate3d(v3.Vec{centerDistance, 0, cfg.ShaftLength/4}))
    
    // For visualization - just show the input gear on shaft
    // (output gear would be on separate shaft in real assembly)
    return sdf.Union3D(shaft, inputGear), nil
}

// CoverAssembly creates the complete cover plate
func CoverAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    body, err := CoverPlate(cfg, ModeBody)
    if err != nil {
        return nil, err
    }
    
    holes, err := CoverPlate(cfg, ModeHole)
    if err != nil {
        return nil, err
    }
    
    return sdf.Difference3D(body, holes), nil
}

//-----------------------------------------------------------------------------
// Main Assembly
//-----------------------------------------------------------------------------

// CompleteJointAssembly assembles all components into final model
func CompleteJointAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // Create all sub-assemblies
    basePlate, err := BasePlateAssembly(cfg)
    if err != nil {
        return nil, err
    }
    
    lowerHousing, err := BearingAssembly(cfg, false)
    if err != nil {
        return nil, err
    }
    
    upperHousing, err := BearingAssembly(cfg, true)
    if err != nil {
        return nil, err
    }
    
    driveTrain, err := DriveTrainAssembly(cfg)
    if err != nil {
        return nil, err
    }
    
    cover, err := CoverAssembly(cfg)
    if err != nil {
        return nil, err
    }
    
    // Position components in assembly
    // Base is at origin (z=0)
    basePlate = sdf.Transform3D(basePlate, sdf.Translate3d(v3.Vec{0, 0, -cfg.BaseThickness / 2}))
    
    // Lower housing sits on base
    lowerHousingHeight := cfg.BearingThickness + cfg.HousingFlangeHeight
    lowerHousing = sdf.Transform3D(lowerHousing, sdf.Translate3d(v3.Vec{0, 0, lowerHousingHeight / 2}))
    
    // Drive train in center
    driveTrain = sdf.Transform3D(driveTrain, sdf.Translate3d(v3.Vec{0, 0, cfg.BearingThickness / 2}))
    
    // Upper housing above drive train
    upperHousing = sdf.Transform3D(upperHousing, sdf.Translate3d(v3.Vec{0, 0, cfg.ShaftLength - lowerHousingHeight/2}))
    
    // Cover on top
    plateThickness := cfg.BaseThickness * 0.6
    cover = sdf.Transform3D(cover, sdf.Translate3d(v3.Vec{0, 0, cfg.ShaftLength + plateThickness/2}))
    
    // Combine all components
    assembly := sdf.Union3D(
        basePlate,
        lowerHousing,
        driveTrain,
        upperHousing,
        cover,
    )
    
    return assembly, nil
}

//-----------------------------------------------------------------------------
// Rendering Functions
//-----------------------------------------------------------------------------

// RenderComponent renders a single component to STL
func RenderComponent(component sdf.SDF3, filename string, cfg *JointConfig) {
    scaled := sdf.ScaleUniform3D(component, cfg.Material.ShrinkageFactor)
    render.ToSTL(scaled, filename, render.NewMarchingCubesOctree(cfg.MeshResolution))
}

// ExportIndividualComponents exports each part separately for manufacturing
func ExportIndividualComponents(cfg *JointConfig) error {
    log.Println("Exporting individual components...")
    
    components := map[string]func(*JointConfig) (sdf.SDF3, error){
        "base_plate":      BasePlateAssembly,
        "lower_housing":   func(c *JointConfig) (sdf.SDF3, error) { return BearingAssembly(c, false) },
        "upper_housing":   func(c *JointConfig) (sdf.SDF3, error) { return BearingAssembly(c, true) },
        "drive_train":     DriveTrainAssembly,
        "cover_plate":     CoverAssembly,
    }
    
    for name, buildFunc := range components {
        log.Printf("Building %s...", name)
        component, err := buildFunc(cfg)
        if err != nil {
            return err
        }
        
        filename := name + ".stl"
        log.Printf("Rendering %s...", filename)
        RenderComponent(component, filename, cfg)
    }
    
    return nil
}

//-----------------------------------------------------------------------------
// Main Entry Point
//-----------------------------------------------------------------------------

func main() {
    log.Println("=== Complex 3D CAD Architecture Demo ===")
    log.Println("Parametric Robotic Joint Assembly")
    log.Println()
    
    // Create configuration
    cfg := NewDefaultJointConfig()
    
    log.Printf("Material: %s (shrinkage: %.3f%%)\n", 
        cfg.Material.Name, 
        (cfg.Material.ShrinkageFactor-1)*100)
    log.Printf("Gear Ratio: %d:%d (%.2f:1)\n", 
        cfg.InputTeeth, 
        cfg.OutputTeeth, 
        float64(cfg.OutputTeeth)/float64(cfg.InputTeeth))
    log.Println()
    
    // Option 1: Export complete assembly for visualization
    log.Println("Building complete assembly...")
    assembly, err := CompleteJointAssembly(cfg)
    if err != nil {
        log.Fatalf("Failed to build assembly: %s", err)
    }
    
    log.Println("Rendering complete assembly...")
    RenderComponent(assembly, "joint_assembly.stl", cfg)
    
    // Option 2: Export individual components for 3D printing
    if err := ExportIndividualComponents(cfg); err != nil {
        log.Fatalf("Failed to export components: %s", err)
    }
    
    log.Println()
    log.Println("✓ Export complete!")
    log.Println("Generated files:")
    log.Println("  - joint_assembly.stl (complete assembly)")
    log.Println("  - base_plate.stl")
    log.Println("  - lower_housing.stl")
    log.Println("  - upper_housing.stl")
    log.Println("  - drive_train.stl")
    log.Println("  - cover_plate.stl")
}

//-----------------------------------------------------------------------------
