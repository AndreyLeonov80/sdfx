//-----------------------------------------------------------------------------
/*

Usage Examples for Complex Architecture

This file shows practical examples of how to use the architectural patterns.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// Example 1: Using Builder Pattern
//-----------------------------------------------------------------------------

// ExampleBuilderPattern demonstrates fluent configuration
func ExampleBuilderPattern() (sdf.SDF3, error) {
	fmt.Println("\n=== Example 1: Builder Pattern ===")
	
	// Create a custom joint using builder pattern
	joint, err := NewJointBuilder().
		WithMaterial("ABS").
		WithBaseDimensions(100, 10).
		WithGearRatio(20, 60).
		WithShaft(15, 60).
		WithQuality(200).
		Build()
	
	if err != nil {
		return nil, err
	}
	
	// Inspect the configuration
	cfg := NewJointBuilder().
		WithMaterial("ABS").
		WithBaseDimensions(100, 10).
		Config()
	
	fmt.Printf("Configuration: Base=%.1fmm, Material=%s, Ratio=%d:%d\n",
		cfg.BaseDiameter,
		cfg.Material.Name,
		cfg.InputTeeth,
		cfg.OutputTeeth)
	
	return joint, nil
}

//-----------------------------------------------------------------------------
// Example 2: Using Component Registry
//-----------------------------------------------------------------------------

// ExampleComponentRegistry demonstrates dynamic component management
func ExampleComponentRegistry() error {
	fmt.Println("\n=== Example 2: Component Registry ===")
	
	registry := NewComponentRegistry()
	cfg := NewDefaultJointConfig()
	
	// List all available components
	fmt.Println("Available components:")
	for _, name := range registry.List() {
		fmt.Printf("  - %s\n", name)
	}
	
	// Build individual components
	fmt.Println("\nBuilding components...")
	for _, name := range []string{"base_plate", "drive_train"} {
		component, err := registry.Build(name, cfg)
		if err != nil {
			return err
		}
		fmt.Printf("  ✓ Built %s\n", name)
		_ = component // Use component
	}
	
	// Register a custom component
	registry.Register("custom_bracket", func(cfg *JointConfig) (sdf.SDF3, error) {
		return sdf.Box3D(v3.Vec{20, 20, 5}, 1)
	})
	
	fmt.Println("\nCustom component registered!")
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 3: Using Assembly Pipeline
//-----------------------------------------------------------------------------

// ExampleAssemblyPipeline demonstrates structured assembly
func ExampleAssemblyPipeline() (sdf.SDF3, error) {
	fmt.Println("\n=== Example 3: Assembly Pipeline ===")
	
	cfg := NewDefaultJointConfig()
	pipeline := NewAssemblyPipeline(cfg)
	
	// Define assembly steps with transforms
	baseTransform := sdf.Translate3d(v3.Vec{0, 0, -cfg.BaseThickness / 2})
	housingTransform := sdf.Translate3d(v3.Vec{0, 0, 10})
	
	pipeline.
		AddStep(&AssemblyStep{
			Name:      "base_plate",
			Build:     BasePlateAssembly,
			Transform: &baseTransform,
		}).
		AddStep(&AssemblyStep{
			Name: "lower_housing",
			Build: func(c *JointConfig) (sdf.SDF3, error) {
				return BearingAssembly(c, false)
			},
			Transform: &housingTransform,
		})
	
	fmt.Println("Pipeline configured with 2 steps")
	
	// Execute pipeline
	assembly, err := pipeline.Execute()
	if err != nil {
		return nil, err
	}
	
	fmt.Println("✓ Pipeline executed successfully")
	
	return assembly, nil
}

//-----------------------------------------------------------------------------
// Example 4: Using Variant Generator
//-----------------------------------------------------------------------------

// ExampleVariantGenerator demonstrates design exploration
func ExampleVariantGenerator() error {
	fmt.Println("\n=== Example 4: Variant Generator ===")
	
	baseConfig := NewDefaultJointConfig()
	generator := NewVariantGenerator(baseConfig)
	
	// Add design variations
	generator.
		AddVariation(ScaleVariant(0.8)).
		AddVariation(ScaleVariant(1.2)).
		AddVariation(MaterialVariant("ABS")).
		AddVariation(GearRatioVariant(15, 45)).
		AddVariation(GearRatioVariant(25, 75))
	
	variants, err := generator.Generate()
	if err != nil {
		return err
	}
	
	fmt.Printf("Generated %d design variants:\n", len(variants))
	for i, variant := range variants {
		fmt.Printf("  Variant %d: Base=%.1fmm, Material=%s, Ratio=%d:%d\n",
			i+1,
			variant.BaseDiameter,
			variant.Material.Name,
			variant.InputTeeth,
			variant.OutputTeeth)
	}
	
	// Could render each variant:
	// for i, variant := range variants {
	//     assembly, _ := CompleteJointAssembly(variant)
	//     filename := fmt.Sprintf("variant_%d.stl", i)
	//     RenderComponent(assembly, filename, variant)
	// }
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 5: Using Component Caching
//-----------------------------------------------------------------------------

// ExampleComponentCaching demonstrates performance optimization
func ExampleComponentCaching() error {
	fmt.Println("\n=== Example 5: Component Caching ===")
	
	cache := NewComponentCache()
	cfg := NewDefaultJointConfig()
	
	// First build - will be cached
	fmt.Println("First build (no cache)...")
	component1, err := cache.GetOrBuild("base_plate", func() (sdf.SDF3, error) {
		fmt.Println("  Building base plate...")
		return BasePlateAssembly(cfg)
	})
	if err != nil {
		return err
	}
	_ = component1
	
	// Second build - retrieved from cache
	fmt.Println("Second build (from cache)...")
	component2, err := cache.GetOrBuild("base_plate", func() (sdf.SDF3, error) {
		fmt.Println("  Building base plate...")
		return BasePlateAssembly(cfg)
	})
	if err != nil {
		return err
	}
	_ = component2
	
	fmt.Printf("Cache size: %d components\n", cache.Size())
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 6: Using Constraints Validation
//-----------------------------------------------------------------------------

// ExampleConstraintsValidation demonstrates design validation
func ExampleConstraintsValidation() error {
	fmt.Println("\n=== Example 6: Constraints Validation ===")
	
	validator := NewConstraintValidator()
	
	// Valid configuration
	validCfg := NewDefaultJointConfig()
	fmt.Println("Validating valid configuration...")
	if errors := validator.Validate(validCfg); len(errors) > 0 {
		fmt.Println("  ✗ Validation failed:")
		for _, err := range errors {
			fmt.Printf("    - %s\n", err)
		}
	} else {
		fmt.Println("  ✓ Configuration is valid")
	}
	
	// Invalid configuration
	invalidCfg := NewDefaultJointConfig()
	invalidCfg.HousingWallThickness = 0.5 // Too thin
	invalidCfg.InputTeeth = 5             // Too few teeth
	
	fmt.Println("\nValidating invalid configuration...")
	if errors := validator.Validate(invalidCfg); len(errors) > 0 {
		fmt.Printf("  ✗ Found %d validation errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("    - %s\n", err)
		}
	} else {
		fmt.Println("  ✓ Configuration is valid")
	}
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 7: Using Feature Flags
//-----------------------------------------------------------------------------

// ExampleFeatureFlags demonstrates conditional assembly
func ExampleFeatureFlags() error {
	fmt.Println("\n=== Example 7: Feature Flags ===")
	
	cfg := NewDefaultJointConfig()
	cfg.MeshResolution = 150 // Lower for speed
	
	// Minimal configuration (base only)
	minimalFlags := &FeatureFlags{
		IncludeCover:         false,
		IncludeMountingHoles: true,
		IncludeKeyways:       false,
		IncludeVentilation:   false,
		IncludeBearings:      true,
		IncludeGears:         false,
		SimplifiedGeometry:   true,
	}
	
	fmt.Println("Building minimal assembly...")
	minimal, err := ConditionalJointAssembly(cfg, minimalFlags)
	if err != nil {
		return err
	}
	
	render.ToSTL(
		sdf.ScaleUniform3D(minimal, cfg.Material.ShrinkageFactor),
		"joint_minimal.stl",
		render.NewMarchingCubesOctree(cfg.MeshResolution),
	)
	fmt.Println("  ✓ Generated joint_minimal.stl")
	
	// Full configuration
	fullFlags := DefaultFeatureFlags()
	
	fmt.Println("Building full assembly...")
	full, err := ConditionalJointAssembly(cfg, fullFlags)
	if err != nil {
		return err
	}
	
	render.ToSTL(
		sdf.ScaleUniform3D(full, cfg.Material.ShrinkageFactor),
		"joint_full.stl",
		render.NewMarchingCubesOctree(cfg.MeshResolution),
	)
	fmt.Println("  ✓ Generated joint_full.stl")
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 8: Creating Custom Components
//-----------------------------------------------------------------------------

// CustomMountBracket creates a custom component integrated with the system
func CustomMountBracket(cfg *JointConfig) (sdf.SDF3, error) {
	// Create a mounting bracket that attaches to the base
	bracketLength := cfg.BaseDiameter * 0.4
	bracketWidth := cfg.BaseThickness * 2
	bracketHeight := cfg.BaseThickness
	
	bracket, err := sdf.Box3D(
		v3.Vec{bracketLength, bracketWidth, bracketHeight},
		cfg.Material.GeneralRounding,
	)
	if err != nil {
		return nil, err
	}
	
	// Add mounting hole
	hole, err := sdf.Cylinder3D(bracketHeight*1.5, cfg.BoltDiameter/2+cfg.BoltClearance, 0)
	if err != nil {
		return nil, err
	}
	
	hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{bracketLength / 3, 0, 0}))
	
	return sdf.Difference3D(bracket, hole), nil
}

// ExampleCustomComponent demonstrates integrating custom components
func ExampleCustomComponent() error {
	fmt.Println("\n=== Example 8: Custom Component Integration ===")
	
	cfg := NewDefaultJointConfig()
	cfg.MeshResolution = 150
	
	// Build base assembly
	fmt.Println("Building base assembly...")
	base, err := BasePlateAssembly(cfg)
	if err != nil {
		return err
	}
	
	// Add custom bracket
	fmt.Println("Adding custom mount bracket...")
	bracket, err := CustomMountBracket(cfg)
	if err != nil {
		return err
	}
	
	// Position bracket on side of base
	offset := v3.Vec{cfg.BaseDiameter / 2, 0, 0}
	bracket = sdf.Transform3D(bracket, sdf.Translate3d(offset))
	
	// Combine
	assembly := sdf.Union3D(base, bracket)
	
	render.ToSTL(
		sdf.ScaleUniform3D(assembly, cfg.Material.ShrinkageFactor),
		"joint_with_bracket.stl",
		render.NewMarchingCubesOctree(cfg.MeshResolution),
	)
	
	fmt.Println("  ✓ Generated joint_with_bracket.stl")
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 9: Material Comparison
//-----------------------------------------------------------------------------

// ExampleMaterialComparison demonstrates material-specific variants
func ExampleMaterialComparison() error {
	fmt.Println("\n=== Example 9: Material Comparison ===")
	
	materials := []string{"PLA", "ABS", "PETG"}
	
	fmt.Println("Comparing materials:")
	for _, matName := range materials {
		mat := StandardMaterials[matName]
		cfg := NewDefaultJointConfig()
		cfg.Material = mat
		
		// Adjust wall thickness based on material requirements
		cfg.HousingWallThickness = mat.MinWallThickness * 2
		
		shrinkPercent := (mat.ShrinkageFactor - 1) * 100
		
		fmt.Printf("\n%s:\n", matName)
		fmt.Printf("  Shrinkage: %.2f%%\n", shrinkPercent)
		fmt.Printf("  Min wall: %.2fmm\n", mat.MinWallThickness)
		fmt.Printf("  Rounding: %.2fmm\n", mat.GeneralRounding)
		fmt.Printf("  Housing wall: %.2fmm\n", cfg.HousingWallThickness)
		
		// Could build and render each:
		// assembly, _ := CompleteJointAssembly(cfg)
		// filename := fmt.Sprintf("joint_%s.stl", matName)
		// RenderComponent(assembly, filename, cfg)
	}
	
	return nil
}

//-----------------------------------------------------------------------------
// Example 10: Parametric Size Series
//-----------------------------------------------------------------------------

// ExampleParametricSeries demonstrates creating a product line
func ExampleParametricSeries() error {
	fmt.Println("\n=== Example 10: Parametric Size Series ===")
	
	// Define size series: small, medium, large
	sizes := []struct {
		name  string
		scale float64
		teeth [2]int // input, output
	}{
		{"small", 0.7, [2]int{15, 30}},
		{"medium", 1.0, [2]int{20, 40}},
		{"large", 1.3, [2]int{25, 50}},
	}
	
	fmt.Println("Product line configurations:")
	
	for _, size := range sizes {
		cfg := NewDefaultJointConfig()
		
		// Scale dimensions
		cfg.BaseDiameter *= size.scale
		cfg.BaseThickness *= size.scale
		cfg.ShaftDiameter *= size.scale
		cfg.ShaftLength *= size.scale
		cfg.BearingOD *= size.scale
		cfg.BearingID *= size.scale
		cfg.BearingThickness *= size.scale
		cfg.HousingWallThickness *= size.scale
		
		// Set gear ratio
		cfg.InputTeeth = size.teeth[0]
		cfg.OutputTeeth = size.teeth[1]
		ratio := float64(size.teeth[1]) / float64(size.teeth[0])
		
		fmt.Printf("\n%s:\n", size.name)
		fmt.Printf("  Base: %.1fmm\n", cfg.BaseDiameter)
		fmt.Printf("  Shaft: %.1fmm × %.1fmm\n", cfg.ShaftDiameter, cfg.ShaftLength)
		fmt.Printf("  Gear ratio: %d:%d (%.1f:1)\n", size.teeth[0], size.teeth[1], ratio)
		
		// Validate configuration
		validator := NewConstraintValidator()
		if errors := validator.Validate(cfg); len(errors) > 0 {
			fmt.Printf("  ⚠ Validation warnings:\n")
			for _, err := range errors {
				fmt.Printf("    - %s\n", err)
			}
		} else {
			fmt.Printf("  ✓ Valid configuration\n")
		}
		
		// Could render:
		// assembly, _ := CompleteJointAssembly(cfg)
		// filename := fmt.Sprintf("joint_%s.stl", size.name)
		// RenderComponent(assembly, filename, cfg)
	}
	
	return nil
}

//-----------------------------------------------------------------------------
// Run all examples
//-----------------------------------------------------------------------------

// RunAllExamples executes all example demonstrations
func RunAllExamples() error {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║  Complex Architecture - Usage Examples        ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	
	examples := []struct {
		name string
		fn   func() error
	}{
		{"Component Registry", ExampleComponentRegistry},
		{"Variant Generator", ExampleVariantGenerator},
		{"Component Caching", ExampleComponentCaching},
		{"Constraints Validation", ExampleConstraintsValidation},
		{"Material Comparison", ExampleMaterialComparison},
		{"Parametric Series", ExampleParametricSeries},
	}
	
	for _, ex := range examples {
		if err := ex.fn(); err != nil {
			return fmt.Errorf("%s failed: %w", ex.name, err)
		}
	}
	
	// Examples that generate files
	fileExamples := []struct {
		name string
		fn   func() error
	}{
		{"Feature Flags", ExampleFeatureFlags},
		{"Custom Component", ExampleCustomComponent},
	}
	
	fmt.Println("\n\nGenerating demonstration files...")
	for _, ex := range fileExamples {
		fmt.Printf("\n%s...\n", ex.name)
		if err := ex.fn(); err != nil {
			return fmt.Errorf("%s failed: %w", ex.name, err)
		}
	}
	
	fmt.Println("\n✓ All examples completed successfully!")
	
	return nil
}

//-----------------------------------------------------------------------------
