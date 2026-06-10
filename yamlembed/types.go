package yamlembed

import (
	"strings"

	_ "gopkg.in/yaml.v2"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f,omitempty"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// type indirection — чтобы не вызвать рекурсию
	type BarAlias Bar
	var alias BarAlias
	if err := unmarshal(&alias); err != nil {
		return err
	}
	*b = Bar(alias)
	b.UpperB = strings.ToUpper(b.B)
	return nil
}

func (b Bar) MarshalYAML() (interface{}, error) {
	// Создаем анонимную структуру-обертку
	return struct {
		B string `yaml:"b"`
		// Добавляем тег flow, чтобы массив записался как [..., ...]
		F []any `yaml:"f,omitempty,flow"`
	}{
		B: b.B,
		F: b.F, // Просто передаем оригинальный слайс как есть
	}, nil
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

// Распаковываем данные отдельно в Foo и в Bar
func (baz *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Передаем весь YAML-узел сначала в Foo
	if err := unmarshal(&baz.Foo); err != nil {
		return err
	}
	// А затем тот же самый узел передаем в Bar (где сработает твой кастомный метод Bar.UnmarshalYAML)
	if err := unmarshal(&baz.Bar); err != nil {
		return err
	}
	return nil
}

// Собираем поля из Foo и Bar вручную, так как кастомный маршалер Bar тоже сломает inline
func (baz Baz) MarshalYAML() (interface{}, error) {
	return struct {
		Foo `yaml:",inline"` // Foo не имеет кастомных методов, ее можно инлайнить
		B   string           `yaml:"b"` // Поля из Bar собираем руками
		OI  []string         `yaml:"oi,omitempty"`
		F   []any            `yaml:"f,omitempty,flow"` // Сохраняем логику форматирования
	}{
		Foo: baz.Foo,
		B:   baz.B,
		OI:  baz.OI,
		F:   baz.F,
	}, nil
}
