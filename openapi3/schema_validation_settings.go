package openapi3

import (
	"sync"
)

// SchemaValidationOption describes options a user has when validating request / response bodies.
type SchemaValidationOption func(*schemaValidationSettings)

type schemaValidationSettings struct {
	failfast                         bool
	multiError                       bool
	asreq, asrep                     bool // exclusive (XOR) fields
	formatValidationEnabled          bool
	unknownPropertyValidationEnabled bool
	patternValidationDisabled        bool
	readOnlyValidationDisabled       bool
	writeOnlyValidationDisabled      bool

	onceSettingDefaults sync.Once
	defaultsSet         func()

	customizeMessageError  func(err *SchemaError) string
	customizeSchemaResolve func(string) *Schema
}

func (s *schemaValidationSettings) schemaResolve(schemaRef *SchemaRef) (value *Schema) {
	if value = schemaRef.Value; value != nil || s.customizeSchemaResolve == nil {
		return
	}

	if ref := schemaRef.Ref; len(ref) > 0 {
		value = s.customizeSchemaResolve(ref)
	}
	return
}

// SetCustomizeSchemaResolve allows to fetch custom schema for not empty ref.
// If the passed function returns *Schema.
func SetCustomizeSchemaResolve(resolve func(string) *Schema) SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.customizeSchemaResolve = resolve }
}

// FailFast returns schema validation errors quicker.
func FailFast() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.failfast = true }
}

func MultiErrors() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.multiError = true }
}

func VisitAsRequest() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.asreq, s.asrep = true, false }
}

func VisitAsResponse() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.asreq, s.asrep = false, true }
}

// EnableUnknownPropertyValidation setting makes Validate return an error when validating documents that miss property schema.
func EnableUnknownPropertyValidation() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.unknownPropertyValidationEnabled = true }
}

// EnableFormatValidation setting makes Validate not return an error when validating documents that mention schema formats that are not defined by the OpenAPIv3 specification.
func EnableFormatValidation() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.formatValidationEnabled = true }
}

// DisablePatternValidation setting makes Validate not return an error when validating patterns that are not supported by the Go regexp engine.
func DisablePatternValidation() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.patternValidationDisabled = true }
}

// DisableReadOnlyValidation setting makes Validate not return an error when validating properties marked as read-only
func DisableReadOnlyValidation() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.readOnlyValidationDisabled = true }
}

// DisableWriteOnlyValidation setting makes Validate not return an error when validating properties marked as write-only
func DisableWriteOnlyValidation() SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.writeOnlyValidationDisabled = true }
}

// DefaultsSet executes the given callback (once) IFF schema validation set default values.
func DefaultsSet(f func()) SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.defaultsSet = f }
}

// SetSchemaErrorMessageCustomizer allows to override the schema error message.
// If the passed function returns an empty string, it returns to the previous Error() implementation.
func SetSchemaErrorMessageCustomizer(f func(err *SchemaError) string) SchemaValidationOption {
	return func(s *schemaValidationSettings) { s.customizeMessageError = f }
}

func newSchemaValidationSettings(opts ...SchemaValidationOption) *schemaValidationSettings {
	settings := &schemaValidationSettings{}
	for _, opt := range opts {
		opt(settings)
	}
	return settings
}
