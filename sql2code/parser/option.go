package parser

type NullStyle int

const (
	NullDisable NullStyle = iota
	NullInSql
	NullInPointer
)

type Option func(*options)

type options struct {
	Charset        string
	Collation      string
	JsonTag        bool
	JsonNamedType  int // json命名类型，0:默认，其他值表示驼峰
	TablePrefix    string
	ColumnPrefix   string
	NoNullType     bool
	NullStyle      NullStyle
	Package        string
	GormType       bool
	ForceTableName bool
	IsEmbed        bool // 是否嵌入gorm.Model
}

var defaultOptions = options{
	NullStyle: NullInSql,
	Package:   "model",
}

func WithCharset(charset string) Option {
	return func(o *options) {
		o.Charset = charset
	}
}

func WithCollation(collation string) Option {
	return func(o *options) {
		o.Collation = collation
	}
}

func WithTablePrefix(p string) Option {
	return func(o *options) {
		o.TablePrefix = p
	}
}

func WithColumnPrefix(p string) Option {
	return func(o *options) {
		o.ColumnPrefix = p
	}
}

// WithJsonTag json名称命名类型，0:表示默认，其他值表示驼峰
func WithJsonTag(namedType int) Option {
	return func(o *options) {
		o.JsonTag = true
		o.JsonNamedType = namedType
	}
}

func WithNoNullType() Option {
	return func(o *options) {
		o.NoNullType = true
	}
}

func WithNullStyle(s NullStyle) Option {
	return func(o *options) {
		o.NullStyle = s
	}
}

func WithPackage(pkg string) Option {
	return func(o *options) {
		o.Package = pkg
	}
}

// WithGormType will write type in gorm tag
func WithGormType() Option {
	return func(o *options) {
		o.GormType = true
	}
}

func WithForceTableName() Option {
	return func(o *options) {
		o.ForceTableName = true
	}
}

// WithEmbed is embed gorm.Model
func WithEmbed() Option {
	return func(o *options) {
		o.IsEmbed = true
	}
}

func parseOption(options []Option) options {
	o := defaultOptions
	for _, f := range options {
		f(&o)
	}
	if o.NoNullType {
		o.NullStyle = NullDisable
	}
	return o
}
