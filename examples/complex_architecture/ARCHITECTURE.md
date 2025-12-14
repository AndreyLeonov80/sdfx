# Архитектура Сложных 3D CAD Моделей / Complex 3D CAD Architecture

## Диаграмма Архитектуры / Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         USER APPLICATION                             │
│                          (main.go)                                   │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
        ┌───────────▼──────────┐   ┌─────────▼──────────┐
        │   Configuration      │   │   Rendering        │
        │   Management         │   │   System           │
        │                      │   │                    │
        │ • JointConfig        │   │ • RenderComponent  │
        │ • MaterialConfig     │   │ • Export*          │
        │ • StandardMaterials  │   │                    │
        └──────────┬───────────┘   └────────────────────┘
                   │
                   │
        ┌──────────▼────────────────────────────────────┐
        │         DESIGN PATTERNS                       │
        │  (advanced_patterns.go)                       │
        │                                               │
        │  ┌─────────────┐  ┌─────────────┐           │
        │  │   Builder   │  │  Registry   │           │
        │  │   Pattern   │  │   Pattern   │           │
        │  └─────────────┘  └─────────────┘           │
        │                                              │
        │  ┌─────────────┐  ┌─────────────┐          │
        │  │  Pipeline   │  │  Variant    │          │
        │  │   Pattern   │  │ Generator   │          │
        │  └─────────────┘  └─────────────┘          │
        │                                             │
        │  ┌─────────────┐  ┌─────────────┐         │
        │  │   Caching   │  │ Constraints │         │
        │  │   System    │  │ Validation  │         │
        │  └─────────────┘  └─────────────┘         │
        └──────────┬────────────────────────────────┘
                   │
                   │
        ┌──────────▼───────────────────────────────────┐
        │      COMPONENT LIBRARY                       │
        │                                              │
        │  ┌──────────────┐  ┌───────────────┐       │
        │  │  BaseMount   │  │ BearingHousing│       │
        │  │  • Body      │  │ • Body         │       │
        │  │  • Holes     │  │ • Holes        │       │
        │  └──────────────┘  └───────────────┘       │
        │                                             │
        │  ┌──────────────┐  ┌───────────────┐      │
        │  │  DriveShaft  │  │ParametricGear │      │
        │  │  • Keyway    │  │ • Teeth        │      │
        │  └──────────────┘  └───────────────┘      │
        │                                            │
        │  ┌──────────────┐  ┌───────────────┐     │
        │  │  CoverPlate  │  │CounterboreHole│     │
        │  │  • Vents     │  │ • Bolt Head    │     │
        │  └──────────────┘  └───────────────┘     │
        └──────────┬─────────────────────────────────┘
                   │
                   │
        ┌──────────▼──────────────────────────────────┐
        │       SUB-ASSEMBLIES                        │
        │                                             │
        │  • BasePlateAssembly()                     │
        │  • BearingAssembly()                       │
        │  • DriveTrainAssembly()                    │
        │  • CoverAssembly()                         │
        └──────────┬────────────────────────────────┘
                   │
                   │
        ┌──────────▼──────────────────────────────────┐
        │     MAIN ASSEMBLY                           │
        │                                             │
        │  CompleteJointAssembly()                   │
        │  • Position components                     │
        │  • Apply transforms                        │
        │  • Union operations                        │
        └──────────┬────────────────────────────────┘
                   │
                   │
        ┌──────────▼──────────────────────────────────┐
        │         SDFX LIBRARY                        │
        │                                             │
        │  • sdf.SDF2 / sdf.SDF3                     │
        │  • Boolean operations                      │
        │  • Transformations                         │
        │  • Primitives                              │
        │  • Rendering (Marching Cubes)              │
        └─────────────────────────────────────────────┘
```

## Поток Данных / Data Flow

```
Configuration → Component Library → Sub-Assemblies → Main Assembly → Render → STL
      ↓                ↓                  ↓               ↓            ↓
  Material         Mode Logic        Positioning     Validation    Export
  Settings         (Body/Hole)       Transforms      Constraints   Files
```

## Слои Архитектуры / Architecture Layers

### Слой 1: Конфигурация (Configuration Layer)
- **Назначение**: Управление всеми параметрами модели
- **Компоненты**: `JointConfig`, `MaterialConfig`, константы
- **Принцип**: Single source of truth для всех размеров

### Слой 2: Компоненты (Component Library Layer)
- **Назначение**: Переиспользуемые механические части
- **Паттерн**: Mode-based generation (Body, Hole, Boss, Clearance)
- **Принцип**: Каждый компонент - чистая функция

### Слой 3: Подсборки (Sub-Assembly Layer)
- **Назначение**: Логическая группировка компонентов
- **Операции**: Difference, Union с плавными переходами
- **Принцип**: Модульная композиция

### Слой 4: Главная Сборка (Main Assembly Layer)
- **Назначение**: Финальная интеграция всех частей
- **Операции**: Позиционирование, трансформации, объединение
- **Принцип**: Координация подсборок

### Слой 5: Рендеринг (Rendering Layer)
- **Назначение**: Генерация STL/3MF файлов
- **Технология**: Marching Cubes / Dual Contouring
- **Принцип**: Независимость от логики модели

## Паттерны Проектирования / Design Patterns

### 1. Builder Pattern 🏗️
```go
joint := NewJointBuilder().
    WithMaterial("PLA").
    WithGearRatio(20, 40).
    Build()
```
**Применение**: Удобное создание сложных конфигураций

### 2. Factory Pattern 🏭
```go
cfg := NewDefaultJointConfig()
cfg := NewCustomJointConfig(params...)
```
**Применение**: Создание предустановленных конфигураций

### 3. Registry Pattern 📚
```go
registry.Register("component_name", factoryFunc)
component := registry.Build("component_name", cfg)
```
**Применение**: Динамическое управление компонентами

### 4. Strategy Pattern 🎯
```go
func Component(cfg *Config, mode ComponentMode) SDF3 {
    switch mode {
    case ModeBody: // ...
    case ModeHole: // ...
    }
}
```
**Применение**: Разные стратегии генерации компонента

### 5. Pipeline Pattern ⚙️
```go
pipeline.AddStep(step1).AddStep(step2).Execute()
```
**Применение**: Последовательная сборка с валидацией

### 6. Observer Pattern 👁️
```go
type Constraint interface {
    Validate(cfg *Config) error
}
```
**Применение**: Проверка ограничений проектирования

## Диаграмма Классов / Class Diagram

```
┌─────────────────────┐
│  JointConfig        │
├─────────────────────┤
│ + Material          │
│ + BaseDiameter      │◄─────┐
│ + ShaftDiameter     │      │
│ + GearModule        │      │
│ + ...               │      │
└─────────────────────┘      │
                             │
                             │ uses
┌─────────────────────┐      │
│  Component          │      │
├─────────────────────┤      │
│ + BaseMount()       │──────┘
│ + BearingHousing()  │
│ + DriveShaft()      │
│ + ...               │
└──────────┬──────────┘
           │ composes
           ▼
┌─────────────────────┐
│  SubAssembly        │
├─────────────────────┤
│ + BasePlate()       │
│ + BearingAssembly() │
│ + DriveTrainAsm()   │
└──────────┬──────────┘
           │ combines
           ▼
┌─────────────────────┐
│  MainAssembly       │
├─────────────────────┤
│ + Complete()        │
│ + Export()          │
└─────────────────────┘
```

## Диаграмма Последовательности / Sequence Diagram

```
User          Config      Component    SubAssembly   MainAssembly   Render
 │              │             │            │              │           │
 ├─New Config──►│             │            │              │           │
 │              │             │            │              │           │
 ├─Build Base──────────────► │            │              │           │
 │              │             ├─Create─────►              │           │
 │              │             │            │              │           │
 ├─Build Housing──────────────►            │              │           │
 │              │             ├─Create─────►              │           │
 │              │             │            │              │           │
 ├─Assemble────────────────────────────────────────────► │           │
 │              │             │            │              ├─Position  │
 │              │             │            │              ├─Union     │
 │              │             │            │              │           │
 ├─Export───────────────────────────────────────────────►├─ToSTL────►│
 │              │             │            │              │           ├─Write
 │              │             │            │              │           │
 ◄──────────────────────────────────────────────────────────────────┘
```

## Расширение Системы / System Extension

### Добавление Нового Компонента

```
1. Определить параметры в JointConfig
   ↓
2. Создать функцию компонента
   func NewComponent(cfg *Config, mode Mode) SDF3
   ↓
3. Добавить в подсборку или главную сборку
   ↓
4. Зарегистрировать в реестре (опционально)
   ↓
5. Протестировать и экспортировать
```

### Добавление Нового Паттерна

```
1. Определить интерфейс паттерна
   ↓
2. Реализовать основную логику
   ↓
3. Создать примеры использования
   ↓
4. Добавить в documentation
```

## Лучшие Практики / Best Practices

### ✅ DO (Делать)

1. **Параметризация**: Все размеры через конфигурацию
2. **Модульность**: Компоненты как независимые функции
3. **Валидация**: Проверка ограничений перед сборкой
4. **Именование**: Понятные имена функций и переменных
5. **Документация**: Комментарии для сложной логики
6. **Тестирование**: Unit-тесты для компонентов

### ❌ DON'T (Не делать)

1. **Хардкод**: Жёстко заданные значения в коде
2. **Монолиты**: Гигантские функции на 500+ строк
3. **Глобальные**: Глобальное изменяемое состояние
4. **Дублирование**: Копипаст кода компонентов
5. **Без ошибок**: Игнорирование error handling
6. **Без проверок**: Отсутствие валидации входных данных

## Метрики Качества / Quality Metrics

- **Модульность**: Каждый компонент < 100 строк
- **Переиспользование**: > 50% кода переиспользуется
- **Тестирование**: > 80% покрытие тестами
- **Документация**: Комментарии для всех публичных функций
- **Производительность**: Кэширование часто используемых компонентов
- **Гибкость**: Легко создавать вариации (< 10 строк кода)

## Примеры Использования / Usage Examples

См. файл `examples.go` для 10+ практических примеров применения архитектуры.

## Ссылки / References

- [main.go](main.go) - Основная реализация
- [advanced_patterns.go](advanced_patterns.go) - Паттерны проектирования
- [examples.go](examples.go) - Примеры использования
- [README.md](README.md) - Документация пользователя
