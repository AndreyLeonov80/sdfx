# Complex 3D CAD Architecture Example

## Обзор / Overview

Этот пример демонстрирует профессиональную архитектуру кода для сложных 3D CAD моделей с использованием библиотеки SDFX.

This example demonstrates professional code architecture for complex 3D CAD models using the SDFX library.

## Архитектурные Принципы / Architectural Principles

### 1. **Управление Конфигурацией / Configuration Management**

```go
type JointConfig struct {
    Material      MaterialConfig
    BaseDiameter  float64
    // ... все параметры в одной структуре
}
```

**Преимущества:**
- Централизованное управление параметрами
- Легко создавать вариации модели
- Простое тестирование разных конфигураций
- Возможность сериализации в JSON/YAML

### 2. **Библиотека Компонентов / Component Library**

Каждый компонент - это функция, возвращающая `sdf.SDF3`:

```go
func BaseMount(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error)
func BearingHousing(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error)
func DriveShaft(cfg *JointConfig, includeKeyway bool) (sdf.SDF3, error)
```

**Паттерн Mode:**
- `ModeBody` - основное тело компонента
- `ModeHole` - отверстия и вырезы
- `ModeBoss` - выступы и бобышки
- `ModeClearance` - зазоры для сборки

### 3. **Подсборки / Sub-Assemblies**

Логическая группировка связанных компонентов:

```go
func BasePlateAssembly(cfg *JointConfig) (sdf.SDF3, error)
func BearingAssembly(cfg *JointConfig, upperHousing bool) (sdf.SDF3, error)
func DriveTrainAssembly(cfg *JointConfig) (sdf.SDF3, error)
```

**Преимущества:**
- Модульная структура
- Переиспользование кода
- Независимое тестирование
- Параллельная разработка

### 4. **Главная Сборка / Main Assembly**

```go
func CompleteJointAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // Создание всех подсборок
    basePlate, _ := BasePlateAssembly(cfg)
    lowerHousing, _ := BearingAssembly(cfg, false)
    
    // Позиционирование компонентов
    basePlate = sdf.Transform3D(basePlate, sdf.Translate3d(...))
    
    // Объединение
    return sdf.Union3D(basePlate, lowerHousing, ...), nil
}
```

### 5. **Система Рендеринга / Rendering System**

```go
func RenderComponent(component sdf.SDF3, filename string, cfg *JointConfig)
func ExportIndividualComponents(cfg *JointConfig)
```

## Структура Проекта / Project Structure

```
complex_architecture/
├── main.go           # Основной файл с архитектурой
├── README.md         # Документация (этот файл)
└── [generated STLs]  # Сгенерированные 3D модели
```

## Использование / Usage

### Компиляция и запуск / Build and Run

```bash
cd examples/complex_architecture
go run main.go
```

### Создание вариантов / Creating Variants

```go
func main() {
    // Вариант 1: Маленький сустав
    smallCfg := NewDefaultJointConfig()
    smallCfg.BaseDiameter = 60.0
    smallCfg.InputTeeth = 15
    smallCfg.OutputTeeth = 30
    
    // Вариант 2: Большой сустав из ABS
    largeCfg := NewDefaultJointConfig()
    largeCfg.Material = StandardMaterials["ABS"]
    largeCfg.BaseDiameter = 120.0
    largeCfg.InputTeeth = 30
    largeCfg.OutputTeeth = 60
    
    // Рендеринг обоих вариантов
    assembly1, _ := CompleteJointAssembly(smallCfg)
    assembly2, _ := CompleteJointAssembly(largeCfg)
}
```

## Лучшие Практики / Best Practices

### 1. Параметризация

❌ **Плохо:**
```go
cylinder, _ := sdf.Cylinder3D(50.0, 25.0, 0.5)
```

✅ **Хорошо:**
```go
cylinder, _ := sdf.Cylinder3D(
    cfg.ShaftLength,
    cfg.ShaftDiameter/2,
    cfg.Material.GeneralRounding,
)
```

### 2. Обработка Ошибок

❌ **Плохо:**
```go
component, _ := BuildComponent(cfg)
```

✅ **Хорошо:**
```go
component, err := BuildComponent(cfg)
if err != nil {
    return nil, fmt.Errorf("failed to build component: %w", err)
}
```

### 3. Модульность

❌ **Плохо:**
```go
func CompleteModel() sdf.SDF3 {
    // 500 строк кода создания всей модели
}
```

✅ **Хорошо:**
```go
func CompleteModel() sdf.SDF3 {
    base := BaseAssembly()
    housing := HousingAssembly()
    cover := CoverAssembly()
    return sdf.Union3D(base, housing, cover)
}
```

### 4. Именование

```go
// Константы - UPPER_CASE или camelCase
const BaseMountingHoles = 6

// Конфигурация - структуры с суффиксом Config
type MaterialConfig struct { ... }

// Функции компонентов - существительные
func BearingHousing() sdf.SDF3

// Функции сборок - с суффиксом Assembly
func DriveTrainAssembly() sdf.SDF3
```

### 5. Комментарии

```go
// BaseMount creates the mounting plate for the joint.
// Mode determines what features to generate:
//   - ModeBody: main plate structure
//   - ModeHole: mounting and shaft holes
func BaseMount(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error)
```

## Паттерны Проектирования / Design Patterns

### 1. Factory Pattern (Фабрика)

```go
func NewDefaultJointConfig() *JointConfig { ... }
func NewCustomJointConfig(params ...) *JointConfig { ... }
```

### 2. Builder Pattern (Строитель)

```go
type JointBuilder struct {
    config *JointConfig
}

func (b *JointBuilder) WithMaterial(mat MaterialConfig) *JointBuilder {
    b.config.Material = mat
    return b
}

func (b *JointBuilder) Build() (sdf.SDF3, error) {
    return CompleteJointAssembly(b.config)
}
```

### 3. Strategy Pattern (Стратегия)

```go
type ComponentMode int // Стратегия генерации (Body, Hole, Boss)

func BaseMount(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    switch mode {
    case ModeBody: // ...
    case ModeHole: // ...
    }
}
```

## Расширение Архитектуры / Extending the Architecture

### Добавление Нового Компонента

```go
// 1. Определить параметры в JointConfig
type JointConfig struct {
    // ...
    NewComponentSize float64
}

// 2. Создать функцию компонента
func NewComponent(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    switch mode {
    case ModeBody:
        // Создать тело
    case ModeHole:
        // Создать отверстия
    }
}

// 3. Добавить в сборку
func CompleteJointAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // ...
    newComp, err := NewComponent(cfg, ModeBody)
    // Позиционировать и объединить
}
```

### Добавление Нового Материала

```go
StandardMaterials["Nylon"] = MaterialConfig{
    Name:             "Nylon",
    ShrinkageFactor:  1.0 / 0.985,
    MinWallThickness: 1.8,
    GeneralRounding:  0.6,
}
```

## Производительность / Performance

### Кэширование Компонентов

```go
type CachedComponent struct {
    component sdf.SDF3
    hash      string
}

var componentCache = make(map[string]sdf.SDF3)
```

### Уровни Детализации (LOD)

```go
type QualityLevel int

const (
    QualityDraft QualityLevel = iota  // 100 res
    QualityNormal                      // 300 res
    QualityHigh                        // 500 res
)

func (cfg *JointConfig) SetQuality(level QualityLevel) {
    switch level {
    case QualityDraft:
        cfg.MeshResolution = 100
    case QualityNormal:
        cfg.MeshResolution = 300
    case QualityHigh:
        cfg.MeshResolution = 500
    }
}
```

## Тестирование / Testing

```go
func TestBasePlateAssembly(t *testing.T) {
    cfg := NewDefaultJointConfig()
    
    plate, err := BasePlateAssembly(cfg)
    if err != nil {
        t.Fatalf("Failed to build base plate: %v", err)
    }
    
    // Проверить размеры
    bounds := plate.BoundingBox()
    if bounds.Size().X > cfg.BaseDiameter {
        t.Errorf("Base plate too large")
    }
}
```

## Экспорт / Export

Сгенерированные файлы:

1. **joint_assembly.stl** - Полная сборка для визуализации
2. **base_plate.stl** - Базовая плита
3. **lower_housing.stl** - Нижний корпус подшипника
4. **upper_housing.stl** - Верхний корпус подшипника
5. **drive_train.stl** - Вал с шестерней
6. **cover_plate.stl** - Защитная крышка

## Полезные Ссылки / Useful Links

- [SDFX Documentation](https://github.com/deadsy/sdfx)
- [Signed Distance Functions](https://iquilezles.org/articles/distfunctions/)
- [3D Printing Best Practices](https://www.simplify3d.com/resources/print-quality-troubleshooting/)

## Лицензия / License

Используйте этот код как шаблон для своих проектов.
Use this code as a template for your projects.
