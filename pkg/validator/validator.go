package validator

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"time"

	// "github.com/go-playground/locales/en"
	// ut "github.com/go-playground/universal-translator"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/go-playground/validator/v10"

	// en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	}
	if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
		log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	}
	if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
		log.Fatal().Err(err).Msg("Error while registering unique validator")
	}
	if err := v.RegisterValidation("phone", phoneValidator); err != nil {
		log.Fatal().Err(err).Msg("Error while registering phone validator")
	}
	if err := v.RegisterValidation("otp_number", otpNumberValidation); err != nil {
		log.Fatal().Err(err).Msg("Error while registering otp_number validator")
	}
	if err := v.RegisterValidation("google_token", isGoogleTokenValid); err != nil {
		log.Fatal().Err(err).Msg("Error while registering google_token validator")
	}
	if err := v.RegisterValidation("role_permission_action", isValidRolePermissionAction); err != nil {
		log.Fatal().Err(err).Msg("Error while registering role_permission_action validator")
	}
	if err := v.RegisterValidation("resource_permission_action", isValidResourcePermissionAction); err != nil {
		log.Fatal().Err(err).Msg("Error while registering resource_permission_action validator")
	}
	if err := v.RegisterValidation("device_type", isValidDeviceType); err != nil {
		log.Fatal().Err(err).Msg("Error while registering device_type validator")
	}
	if err := v.RegisterValidation("non_empty_array", isNonEmptyArray); err != nil {
		log.Fatal().Err(err).Msg("Error while registering non_empty_array validator")
	}
	if err := v.RegisterValidation("xss_safe", isXSSSafe); err != nil {
		log.Fatal().Err(err).Msg("Error while registering xss_safe validator")
	}
	if err := v.RegisterValidation("json_string", isValidJSON); err != nil {
		log.Fatal().Err(err).Msg("Error while registering json_string validator")
	}
	if err := v.RegisterValidation("non_zero_integer", isNonZeroInt); err != nil {
		log.Fatal().Err(err).Msg("Error while registering non_zero_integer validator")
	}
	if err := v.RegisterValidation("date_format", isValidTimeFormat); err != nil {
		log.Fatal().Err(err).Msg("Error while registering date_format validator")
	}
	if err := v.RegisterValidation("number", isInteger); err != nil {
		log.Fatal().Err(err).Msg("Error while registering integer validator")
	}
	if err := v.RegisterValidation("callback_finish", isValidCallbackFinish); err != nil {
		log.Fatal().Err(err).Msg("Error while registering callback finish validator")
	}

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

// blacklist email validator
func isEmailBlacklist(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	disallowedDomains := []string{"outlook", "hotmail", "aol", "live", "inbox", "icloud", "mail", "gmx", "yandex"}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.Split(parts[1], ".")[0]

	for _, disallowed := range disallowedDomains {
		if domain == disallowed {
			return false
		}
	}

	return true
}

func isStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= 'a' && char <= 'z':
			hasLowercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUppercase && hasLowercase && hasNumber
}

func isUniqueInSlice(fl validator.FieldLevel) bool {
	// Get the slice from the FieldLevel interface
	val := fl.Field()

	// Ensure the field is a slice
	if val.Kind() != reflect.Slice {
		return false
	}

	// Use a map to check for duplicates
	elements := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if _, found := elements[elem]; found {
			return false // Duplicate found
		}
		elements[elem] = true
	}
	return true
}

func phoneValidator(fl validator.FieldLevel) bool {
	// +62, 62, 0
	phoneRegex := `^(?:\+62|62|0)[2-9][0-9]{7,14}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(fl.Field().String())
}

func otpNumberValidation(fl validator.FieldLevel) bool {
	otp := fl.Field().String()
	matched, _ := regexp.MatchString(`^[0-9A-Z]{6}$`, otp)
	return matched
}

func isGoogleTokenValid(fl validator.FieldLevel) bool {
	token := fl.Field().String()

	if len(token) < 100 || len(token) > 1000 {
		return false
	}

	jwtRegex := `^[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+$`
	matched, err := regexp.MatchString(jwtRegex, token)
	if err != nil || !matched {
		return false
	}

	return true
}

func isValidRolePermissionAction(fl validator.FieldLevel) bool {
	allowedActions := []string{
		"view_all",
		"view_on",
		"create",
		"update",
		"read",
		"delete",
		"use",
	}

	action := fl.Field().String()

	for _, allowed := range allowedActions {
		if action == allowed {
			return true
		}
	}

	return false
}

func isValidResourcePermissionAction(fl validator.FieldLevel) bool {
	allowedResources := []string{
		"notification",
		"read_notification",
		"voucher",
		"users",
		"subcategory",
		"category",
		"product",
		"roles",
		"permissions",
		"banner",
		"expenses",
		"order",
		"dashboard_ecommerce",
		"product_content",
		"article",
		"faq",
		"category_article",
		"dashboard",
		"profile",
		"address",
		"recipe",
		"product_pos",
		"expenses_pos",
		"ingredient",
		"dashboard_pos",
		"product_category_pos",
		"member",
		"transaction_history_pos",
	}

	resource := fl.Field().String()

	for _, allowed := range allowedResources {
		if resource == allowed {
			return true
		}
	}

	return false
}

func isValidDeviceType(fl validator.FieldLevel) bool {
	allowedDeviceType := []string{
		"android",
		"ios",
	}

	deviceType := fl.Field().String()

	for _, allowed := range allowedDeviceType {
		if strings.EqualFold(deviceType, allowed) {
			return true
		}
	}

	return false
}

func isNonEmptyArray(fl validator.FieldLevel) bool {
	val := fl.Field()

	if val.Kind() != reflect.Slice {
		return false
	}

	return val.Len() > 0
}

func isXSSSafe(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	p := bluemonday.UGCPolicy()
	sanitized := p.Sanitize(input)

	return input == sanitized
}

func isValidJSON(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	var js json.RawMessage

	return json.Unmarshal([]byte(input), &js) == nil
}

func isNonZeroInt(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() > 0
	default:
		return false
	}
}

func isInteger(fl validator.FieldLevel) bool {
	kind := fl.Field().Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isValidTimeFormat(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "" {
		return false
	}

	_, err := time.Parse("2006-01-02T15:04:05.000Z", val)
	return err == nil
}

func isValidCallbackFinish(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return val == config.Envs.App.Domain
}
