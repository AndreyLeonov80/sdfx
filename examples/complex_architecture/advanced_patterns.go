//-----------------------------------------------------------------------------
/*

Advanced Architectural Patterns for Complex CAD Models

This file demonstrates additional design patterns and techniques
for organizing large-scale 3D CAD projects.

*/
//-----------------------------------------------------------------------------

package main

import (
    "fmt"
    "sync"

    "github.com/deadsy/sdfx/sdf"
    v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// Pattern 1: Builder Pattern for Complex Configurations
//-----------------------------------------------------------------------------

// JointBuilder provides a fluent interface for building joint configurations
type JointBuilder struct {
    config *JointConfig
}

// NewJointBuilder creates a new builder with default configuration
func NewJointBuilder() *JointBuilder {
    return &JointBuilder{
        config: NewDefaultJointConfig(),
    }
}

// WithMaterial sets the material for the joint
func (b *JointBuilder) WithMaterial(materialName string) *JointBuilder {
    if mat, ok := StandardMaterials[materialName]; ok {
        b.config.Material = mat
    }
    return b
}

// WithBaseDimensions sets base plate dimensions
func (b *JointBuilder) WithBaseDimensions(diameter, thickness float64) *JointBuilder {
    b.config.BaseDiameter = diameter
    b.config.BaseThickness = thickness
    return b
}

// WithGearRatio sets the gear ratio (input:output)
func (b *JointBuilder) WithGearRatio(inputTeeth, outputTeeth int) *JointBuilder {
    b.config.InputTeeth = inputTeeth
    b.config.OutputTeeth = outputTeeth
    return b
}

// WithShaft sets shaft dimensions
func (b *JointBuilder) WithShaft(diameter, length float64) *JointBuilder {
    b.config.ShaftDiameter = diameter
    b.config.ShaftLength = length
    return b
}

// WithQuality sets rendering quality
func (b *JointBuilder) WithQuality(resolution int) *JointBuilder {
    b.config.MeshResolution = resolution
    return b
}

// Build creates the complete assembly
func (b *JointBuilder) Build() (sdf.SDF3, error) {
    return CompleteJointAssembly(b.config)
}

// Config returns the current configuration (useful for inspection/validation)
func (b *JointBuilder) Config() *JointConfig {
    return b.config
}

// Example usage:
// joint, err := NewJointBuilder().
//     WithMaterial("ABS").
//     WithBaseDimensions(100, 10).
//     WithGearRatio(20, 60).
//     WithQuality(400).
//     Build()

//-----------------------------------------------------------------------------
// Pattern 2: Component Registry for Dynamic Component Management
//-----------------------------------------------------------------------------

// ComponentFactory is a function that creates a component
type ComponentFactory func(*JointConfig) (sdf.SDF3, error)

// ComponentRegistry manages available components
type ComponentRegistry struct {
    components map[string]ComponentFactory
    mu         sync.RWMutex
}

// NewComponentRegistry creates a new registry
func NewComponentRegistry() *ComponentRegistry {
    registry := &ComponentRegistry{
        components: make(map[string]ComponentFactory),
    }
    
    // Register default components
    registry.Register("base_plate", BasePlateAssembly)
    registry.Register("lower_housing", func(cfg *JointConfig) (sdf.SDF3, error) {
        return BearingAssembly(cfg, false)
    })
    registry.Register("upper_housing", func(cfg *JointConfig) (sdf.SDF3, error) {
        return BearingAssembly(cfg, true)
    })
    registry.Register("drive_train", DriveTrainAssembly)
    registry.Register("cover_plate", CoverAssembly)
    
    return registry
}

// Register adds a component factory to the registry
func (r *ComponentRegistry) Register(name string, factory ComponentFactory) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.components[name] = factory
}

// Get retrieves a component factory by name
func (r *ComponentRegistry) Get(name string) (ComponentFactory, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    factory, ok := r.components[name]
    return factory, ok
}

// Build creates a component by name
func (r *ComponentRegistry) Build(name string, cfg *JointConfig) (sdf.SDF3, error) {
    factory, ok := r.Get(name)
    if !ok {
        return nil, fmt.Errorf("component '%s' not found in registry", name)
    }
    return factory(cfg)
}

// List returns all registered component names
func (r *ComponentRegistry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    names := make([]string, 0, len(r.components))
    for name := range r.components {
        names = append(names, name)
    }
    return names
}

//-----------------------------------------------------------------------------
// Pattern 3: Assembly Pipeline for Complex Build Sequences
//-----------------------------------------------------------------------------

// AssemblyStep represents a single step in the assembly process
type AssemblyStep struct {
    Name       string
    Build      func(*JointConfig) (sdf.SDF3, error)
    Transform  *sdf.M44  // Pointer to allow nil check
    Validate   func(sdf.SDF3) error
}

// AssemblyPipeline manages the sequence of assembly steps
type AssemblyPipeline struct {
    steps []*AssemblyStep
    cfg   *JointConfig
}

// NewAssemblyPipeline creates a new pipeline
func NewAssemblyPipeline(cfg *JointConfig) *AssemblyPipeline {
    return &AssemblyPipeline{
        steps: make([]*AssemblyStep, 0),
        cfg:   cfg,
    }
}

// AddStep adds a step to the pipeline
func (p *AssemblyPipeline) AddStep(step *AssemblyStep) *AssemblyPipeline {
    p.steps = append(p.steps, step)
    return p
}

// Execute runs the pipeline and combines all components
func (p *AssemblyPipeline) Execute() (sdf.SDF3, error) {
    if len(p.steps) == 0 {
        return nil, fmt.Errorf("pipeline has no steps")
    }
    
    components := make([]sdf.SDF3, 0, len(p.steps))
    
    for i, step := range p.steps {
        // Build component
        component, err := step.Build(p.cfg)
        if err != nil {
            return nil, fmt.Errorf("step %d (%s) failed: %w", i, step.Name, err)
        }
        
        // Validate if validator provided
        if step.Validate != nil {
            if err := step.Validate(component); err != nil {
                return nil, fmt.Errorf("step %d (%s) validation failed: %w", i, step.Name, err)
            }
        }
        
        // Apply transform
        if step.Transform != nil {
            component = sdf.Transform3D(component, *step.Transform)
        }
        
        components = append(components, component)
    }
    
    // Combine all components
    return sdf.Union3D(components...), nil
}

// Example usage:
// pipeline := NewAssemblyPipeline(cfg)
// pipeline.
//     AddStep(&AssemblyStep{
//         Name:  "base",
//         Build: BasePlateAssembly,
//         Transform: sdf.Translate3d(v3.Vec{0, 0, -5}),
//     }).
//     AddStep(&AssemblyStep{
//         Name:  "housing",
//         Build: func(c *JointConfig) (sdf.SDF3, error) {
//             return BearingAssembly(c, false)
//         },
//         Transform: sdf.Translate3d(v3.Vec{0, 0, 10}),
//     })
// assembly, err := pipeline.Execute()

//-----------------------------------------------------------------------------
// Pattern 4: Variant Generator for Design Exploration
//-----------------------------------------------------------------------------

// VariantGenerator creates multiple design variations
type VariantGenerator struct {
    baseConfig *JointConfig
    variations []ConfigModifier
}

// ConfigModifier is a function that modifies a configuration
type ConfigModifier func(*JointConfig)

// NewVariantGenerator creates a new variant generator
func NewVariantGenerator(baseConfig *JointConfig) *VariantGenerator {
    return &VariantGenerator{
        baseConfig: baseConfig,
        variations: make([]ConfigModifier, 0),
    }
}

// AddVariation adds a configuration modifier
func (g *VariantGenerator) AddVariation(modifier ConfigModifier) *VariantGenerator {
    g.variations = append(g.variations, modifier)
    return g
}

// Generate creates all variants
func (g *VariantGenerator) Generate() ([]*JointConfig, error) {
    if len(g.variations) == 0 {
        return nil, fmt.Errorf("no variations defined")
    }
    
    configs := make([]*JointConfig, len(g.variations))
    
    for i, modifier := range g.variations {
        // Create a copy of base config
        cfg := *g.baseConfig
        
        // Apply modification
        modifier(&cfg)
        
        configs[i] = &cfg
    }
    
    return configs, nil
}

// Common variant modifiers
func ScaleVariant(scale float64) ConfigModifier {
    return func(cfg *JointConfig) {
        cfg.BaseDiameter *= scale
        cfg.BaseThickness *= scale
        cfg.ShaftDiameter *= scale
        cfg.ShaftLength *= scale
        cfg.BearingOD *= scale
        cfg.BearingID *= scale
        cfg.BearingThickness *= scale
    }
}

func MaterialVariant(materialName string) ConfigModifier {
    return func(cfg *JointConfig) {
        if mat, ok := StandardMaterials[materialName]; ok {
            cfg.Material = mat
        }
    }
}

func GearRatioVariant(inputTeeth, outputTeeth int) ConfigModifier {
    return func(cfg *JointConfig) {
        cfg.InputTeeth = inputTeeth
        cfg.OutputTeeth = outputTeeth
    }
}

// Example usage:
// generator := NewVariantGenerator(NewDefaultJointConfig())
// generator.
//     AddVariation(ScaleVariant(0.8)).
//     AddVariation(ScaleVariant(1.2)).
//     AddVariation(MaterialVariant("ABS")).
//     AddVariation(GearRatioVariant(15, 45))
// variants, _ := generator.Generate()

//-----------------------------------------------------------------------------
// Pattern 5: Component Caching for Performance
//-----------------------------------------------------------------------------

// ComponentCache caches built components to avoid rebuilding
type ComponentCache struct {
    cache map[string]sdf.SDF3
    mu    sync.RWMutex
}

// NewComponentCache creates a new cache
func NewComponentCache() *ComponentCache {
    return &ComponentCache{
        cache: make(map[string]sdf.SDF3),
    }
}

// Get retrieves a cached component
func (c *ComponentCache) Get(key string) (sdf.SDF3, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    component, ok := c.cache[key]
    return component, ok
}

// Set stores a component in cache
func (c *ComponentCache) Set(key string, component sdf.SDF3) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = component
}

// GetOrBuild retrieves from cache or builds if not present
func (c *ComponentCache) GetOrBuild(key string, builder func() (sdf.SDF3, error)) (sdf.SDF3, error) {
    // Try to get from cache
    if component, ok := c.Get(key); ok {
        return component, nil
    }
    
    // Build component
    component, err := builder()
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    c.Set(key, component)
    
    return component, nil
}

// Clear clears the cache
func (c *ComponentCache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache = make(map[string]sdf.SDF3)
}

// Size returns the number of cached components
func (c *ComponentCache) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.cache)
}

//-----------------------------------------------------------------------------
// Pattern 6: Parametric Constraints and Validation
//-----------------------------------------------------------------------------

// Constraint represents a design constraint
type Constraint interface {
    Validate(cfg *JointConfig) error
    Name() string
}

// MinimumWallThicknessConstraint ensures walls are thick enough
type MinimumWallThicknessConstraint struct {
    MinThickness float64
}

func (c *MinimumWallThicknessConstraint) Name() string {
    return "Minimum Wall Thickness"
}

func (c *MinimumWallThicknessConstraint) Validate(cfg *JointConfig) error {
    if cfg.HousingWallThickness < c.MinThickness {
        return fmt.Errorf("housing wall thickness %.2f is below minimum %.2f",
            cfg.HousingWallThickness, c.MinThickness)
    }
    return nil
}

// BearingFitConstraint ensures bearing fits in housing
type BearingFitConstraint struct{}

func (c *BearingFitConstraint) Name() string {
    return "Bearing Fit"
}

func (c *BearingFitConstraint) Validate(cfg *JointConfig) error {
    housingID := cfg.BearingOD + 2*cfg.BearingClearance
    housingOD := housingID + 2*cfg.HousingWallThickness
    
    if housingOD > cfg.BaseDiameter {
        return fmt.Errorf("bearing housing (%.2f) exceeds base diameter (%.2f)",
            housingOD, cfg.BaseDiameter)
    }
    return nil
}

// GearMeshConstraint ensures gears can mesh properly
type GearMeshConstraint struct{}

func (c *GearMeshConstraint) Name() string {
    return "Gear Mesh"
}

func (c *GearMeshConstraint) Validate(cfg *JointConfig) error {
    if cfg.InputTeeth < 10 {
        return fmt.Errorf("input gear teeth (%d) too few, minimum 10", cfg.InputTeeth)
    }
    if cfg.OutputTeeth < 10 {
        return fmt.Errorf("output gear teeth (%d) too few, minimum 10", cfg.OutputTeeth)
    }
    return nil
}

// ConstraintValidator validates all constraints
type ConstraintValidator struct {
    constraints []Constraint
}

// NewConstraintValidator creates a validator with default constraints
func NewConstraintValidator() *ConstraintValidator {
    return &ConstraintValidator{
        constraints: []Constraint{
            &MinimumWallThicknessConstraint{MinThickness: 1.0},
            &BearingFitConstraint{},
            &GearMeshConstraint{},
        },
    }
}

// AddConstraint adds a custom constraint
func (v *ConstraintValidator) AddConstraint(c Constraint) {
    v.constraints = append(v.constraints, c)
}

// Validate checks all constraints
func (v *ConstraintValidator) Validate(cfg *JointConfig) []error {
    var errors []error
    
    for _, constraint := range v.constraints {
        if err := constraint.Validate(cfg); err != nil {
            errors = append(errors, fmt.Errorf("%s: %w", constraint.Name(), err))
        }
    }
    
    return errors
}

// ValidateOrPanic validates and panics if any constraint fails
func (v *ConstraintValidator) ValidateOrPanic(cfg *JointConfig) {
    if errors := v.Validate(cfg); len(errors) > 0 {
        panic(fmt.Sprintf("Configuration validation failed:\n%v", errors))
    }
}

//-----------------------------------------------------------------------------
// Pattern 7: Feature Flags for Conditional Features
//-----------------------------------------------------------------------------

// FeatureFlags controls which features to include
type FeatureFlags struct {
    IncludeCover          bool
    IncludeMountingHoles  bool
    IncludeKeyways        bool
    IncludeVentilation    bool
    IncludeBearings       bool
    IncludeGears          bool
    SimplifiedGeometry    bool
}

// DefaultFeatureFlags returns all features enabled
func DefaultFeatureFlags() *FeatureFlags {
    return &FeatureFlags{
        IncludeCover:         true,
        IncludeMountingHoles: true,
        IncludeKeyways:       true,
        IncludeVentilation:   true,
        IncludeBearings:      true,
        IncludeGears:         true,
        SimplifiedGeometry:   false,
    }
}

// ConditionalJointAssembly builds assembly based on feature flags
func ConditionalJointAssembly(cfg *JointConfig, flags *FeatureFlags) (sdf.SDF3, error) {
    components := make([]sdf.SDF3, 0)
    
    // Always include base
    basePlate, err := BasePlateAssembly(cfg)
    if err != nil {
        return nil, err
    }
    basePlate = sdf.Transform3D(basePlate, sdf.Translate3d(v3.Vec{0, 0, -cfg.BaseThickness / 2}))
    components = append(components, basePlate)
    
    // Conditional bearing housings
    if flags.IncludeBearings {
        lowerHousing, err := BearingAssembly(cfg, false)
        if err != nil {
            return nil, err
        }
        lowerHousingHeight := cfg.BearingThickness + cfg.HousingFlangeHeight
        lowerHousing = sdf.Transform3D(lowerHousing, sdf.Translate3d(v3.Vec{0, 0, lowerHousingHeight / 2}))
        components = append(components, lowerHousing)
        
        upperHousing, err := BearingAssembly(cfg, true)
        if err != nil {
            return nil, err
        }
        upperHousing = sdf.Transform3D(upperHousing, sdf.Translate3d(v3.Vec{0, 0, cfg.ShaftLength - lowerHousingHeight/2}))
        components = append(components, upperHousing)
    }
    
    // Conditional drive train
    if flags.IncludeGears {
        driveTrain, err := DriveTrainAssembly(cfg)
        if err != nil {
            return nil, err
        }
        driveTrain = sdf.Transform3D(driveTrain, sdf.Translate3d(v3.Vec{0, 0, cfg.BearingThickness / 2}))
        components = append(components, driveTrain)
    }
    
    // Conditional cover
    if flags.IncludeCover {
        cover, err := CoverAssembly(cfg)
        if err != nil {
            return nil, err
        }
        plateThickness := cfg.BaseThickness * 0.6
        cover = sdf.Transform3D(cover, sdf.Translate3d(v3.Vec{0, 0, cfg.ShaftLength + plateThickness/2}))
        components = append(components, cover)
    }
    
    return sdf.Union3D(components...), nil
}

//-----------------------------------------------------------------------------
