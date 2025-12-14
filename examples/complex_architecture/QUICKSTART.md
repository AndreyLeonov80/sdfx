# Quick Start Guide / Быстрый Старт

## 🚀 Запуск / Running

### Вариант 1: Сборка Полной Модели (Default)

```bash
cd examples/complex_architecture
go run main.go
```

**Результат:**
- `joint_assembly.stl` - Полная сборка для визуализации
- `base_plate.stl` - Базовая плита
- `lower_housing.stl` - Нижний корпус подшипника
- `upper_housing.stl` - Верхний корпус подшипника
- `drive_train.stl` - Вал с шестерней
- `cover_plate.stl` - Защитная крышка

### Вариант 2: Запуск Примеров (Examples)

```bash
# Переименовать файлы
mv main.go main_assembly.go.backup
mv main_examples.go.disabled main.go

# Запустить примеры
go run .

# Вернуть обратно
mv main.go main_examples.go.disabled
mv main_assembly.go.backup main.go
```

**Результат:**
- Демонстрация всех паттернов проектирования
- `joint_minimal.stl` - Минимальная конфигурация
- `joint_full.stl` - Полная конфигурация
- `joint_with_bracket.stl` - С кастомным кронштейном

### Вариант 3: Компиляция Бинарника

```bash
go build -o cad_generator
./cad_generator
```

## 📁 Структура Проекта

```
complex_architecture/
├── main.go                    # Главный файл - полная сборка
├── advanced_patterns.go       # Паттерны проектирования
├── examples.go                # Примеры использования
├── main_examples.go.disabled  # Альтернативный main для примеров
├── README.md                  # Детальная документация
├── ARCHITECTURE.md            # Архитектурная диаграмма
└── QUICKSTART.md             # Этот файл
```

## ⚡ Быстрые Модификации

### Изменить Материал

```go
// В main.go найти:
cfg := NewDefaultJointConfig()

// Изменить на:
cfg := NewDefaultJointConfig()
cfg.Material = StandardMaterials["ABS"]  // или "PETG"
```

### Изменить Размер

```go
cfg := NewDefaultJointConfig()
cfg.BaseDiameter = 100.0        // увеличить базу
cfg.ShaftDiameter = 15.0        // толще вал
```

### Изменить Передаточное Отношение

```go
cfg := NewDefaultJointConfig()
cfg.InputTeeth = 15             // входная шестерня
cfg.OutputTeeth = 60            // выходная (60/15 = 4:1)
```

### Изменить Качество Рендеринга

```go
cfg := NewDefaultJointConfig()
cfg.MeshResolution = 200        // быстро (draft)
cfg.MeshResolution = 300        // нормально (default)
cfg.MeshResolution = 500        // высокое качество
```

## 🎨 Создание Собственных Компонентов

### Шаг 1: Добавить Параметры

```go
type JointConfig struct {
    // ... существующие поля ...
    
    // Новые параметры
    MyComponentHeight float64
    MyComponentWidth  float64
}
```

### Шаг 2: Создать Функцию Компонента

```go
func MyCustomComponent(cfg *JointConfig, mode ComponentMode) (sdf.SDF3, error) {
    switch mode {
    case ModeBody:
        return sdf.Box3D(v3.Vec{
            cfg.MyComponentWidth,
            cfg.MyComponentWidth,
            cfg.MyComponentHeight,
        }, cfg.Material.GeneralRounding)
        
    case ModeHole:
        // Добавить отверстия если нужны
        return nil, nil
    }
    return nil, nil
}
```

### Шаг 3: Добавить в Сборку

```go
func CompleteJointAssembly(cfg *JointConfig) (sdf.SDF3, error) {
    // ... существующий код ...
    
    myComponent, err := MyCustomComponent(cfg, ModeBody)
    if err != nil {
        return nil, err
    }
    
    // Позиционировать
    myComponent = sdf.Transform3D(myComponent, 
        sdf.Translate3d(v3.Vec{0, 0, 50}))
    
    // Добавить в сборку
    assembly := sdf.Union3D(
        basePlate,
        // ... другие компоненты ...
        myComponent,
    )
    
    return assembly, nil
}
```

## 🔧 Использование Builder Pattern

```go
func main() {
    // Создать конфигурацию через builder
    joint, err := NewJointBuilder().
        WithMaterial("PETG").
        WithBaseDimensions(90, 8).
        WithGearRatio(25, 50).
        WithShaft(14, 55).
        WithQuality(350).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Экспортировать
    cfg := NewJointBuilder().WithMaterial("PETG").Config()
    RenderComponent(joint, "my_joint.stl", cfg)
}
```

## 📊 Создание Серии Размеров

```go
func main() {
    sizes := map[string]float64{
        "small":  0.7,
        "medium": 1.0,
        "large":  1.3,
    }
    
    for name, scale := range sizes {
        cfg := NewDefaultJointConfig()
        
        // Масштабировать все размеры
        cfg.BaseDiameter *= scale
        cfg.BaseThickness *= scale
        cfg.ShaftDiameter *= scale
        // ...
        
        // Собрать и экспортировать
        assembly, _ := CompleteJointAssembly(cfg)
        filename := fmt.Sprintf("joint_%s.stl", name)
        RenderComponent(assembly, filename, cfg)
    }
}
```

## 🎯 Полезные Команды

```bash
# Проверить синтаксис
go build

# Запустить с выводом информации
go run main.go 2>&1 | tee build.log

# Проверить размер бинарника
go build && ls -lh complex_architecture

# Очистить сгенерированные файлы
rm -f *.stl *.3mf

# Просмотреть STL файл (если установлен meshlab)
meshlab joint_assembly.stl &

# Или используйте любой 3D viewer:
# - Blender
# - FreeCAD
# - PrusaSlicer
# - Cura
```

## 📖 Дополнительная Информация

- **README.md** - Подробная документация с примерами
- **ARCHITECTURE.md** - Диаграммы и архитектурные решения
- **examples.go** - 10+ готовых примеров использования
- **advanced_patterns.go** - Продвинутые паттерны проектирования

## ❓ Частые Вопросы (FAQ)

### Q: Как изменить единицы измерения?
A: Все размеры в миллиметрах. Для изменения умножьте все значения на коэффициент.

### Q: Как экспортировать в 3MF?
A: Замените `render.ToSTL` на `render.To3MF`:
```go
render.To3MF(scaled, "model.3mf", render.NewMarchingCubesOctree(300))
```

### Q: Модель генерируется долго?
A: Уменьшите `MeshResolution` до 150-200 для тестирования.

### Q: Как добавить текст на модель?
A: Используйте `sdf.Text2D()` и `sdf.Extrude3D()`:
```go
text, _ := sdf.Text2D("LOGO", "Arial", 10.0)
text3d := sdf.Extrude3D(text, 2.0)
```

### Q: Можно ли сохранить конфигурацию?
A: Да, сериализуйте `JointConfig` в JSON:
```go
data, _ := json.Marshal(cfg)
ioutil.WriteFile("config.json", data, 0644)
```

## 🎓 Обучающие Ресурсы

1. **Начните с main.go** - изучите базовую структуру
2. **Прочитайте README.md** - поймите концепции
3. **Запустите examples.go** - увидьте паттерны в действии
4. **Измените параметры** - поэкспериментируйте
5. **Создайте свой компонент** - примените знания

## 💡 Советы

- Начинайте с низкого разрешения (150) при разработке
- Используйте валидацию для проверки параметров
- Кэшируйте часто используемые компоненты
- Документируйте сложную логику
- Тестируйте компоненты отдельно перед сборкой

## 🚨 Решение Проблем

**Ошибка: "undefined reference"**
```bash
go mod tidy
go clean -cache
go build
```

**Слишком медленный рендеринг**
```go
cfg.MeshResolution = 150  // Уменьшить разрешение
```

**Файл STL слишком большой**
```go
// Уменьшить разрешение или упростить геометрию
cfg.MeshResolution = 200
```

**Ошибки валидации**
```go
validator := NewConstraintValidator()
errors := validator.Validate(cfg)
for _, err := range errors {
    log.Println(err)
}
```

---

**Готовы начать? Запустите:**
```bash
go run main.go
```

**Нужна помощь? Смотрите:**
- README.md
- ARCHITECTURE.md
- examples.go
