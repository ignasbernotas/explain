package man

import (
	"regexp"
	"strings"
)

var paramPattern = regexp.MustCompile(`(?i)<(.*)>`)

type Option struct {
	Name          string
	NameIndicator string

	Alias          string
	AliasIndicator string

	Parameters  []string
	Description string
}

func (o *Option) String() string {
	name := o.Name
	if name == o.Alias {
		name = o.NameIndicator + name
	} else {
		name = o.AliasIndicator + o.Alias + ", " + o.NameIndicator + name
	}

	return name
}

func (o *Option) StringWithArg() string {
	return o.String() + o.StringArg()
}

func (o *Option) StringArg() string {
	str := ""
	if len(o.Parameters) > 0 && len(o.Parameters[0]) > 0 {
		str += " <" + o.Parameters[0] + ">"
	}
	return str
}

func NewOption(toolName, name, description string) *Option {
	return buildOption(toolName, name, description)
}

func buildOption(toolName, name, description string) *Option {
	var args []string
	var param string

	// extract <parameter> and remove it from name
	if paramPattern.MatchString(name) {
		param = strings.Trim(paramPattern.FindString(name), "<>")
		name = paramPattern.ReplaceAllString(name, "")
	} else {
		// <value> Ar <argument>
		split := strings.Split(name, " "+MacroArgument+" ")
		if len(split) > 1 {
			name = split[0]
			param = split[1]
		}
	}

	for _, arg := range strings.Split(name, ",") {
		args = append(args, strings.TrimSpace(arg))
	}

	opt := &Option{}
	if len(param) > 0 {
		opt.Parameters = []string{param}
	} else {
		opt.Parameters = []string{}
	}

	// handle arg having a short and a long name
	if len(args) == 2 {
		if len(args[0]) < len(args[1]) {
			opt.Name = args[1]
			opt.Alias = args[0]
		}

		if len(args[1]) < len(args[0]) {
			opt.Name = args[0]
			opt.Alias = args[1]
		}
	}

	// if only a single arg defined, use it for both names
	if len(args) == 1 {
		opt.Name = args[0]
		opt.Alias = args[0]
	}

	opt.NameIndicator = getIndicator(opt.Name)
	opt.AliasIndicator = getIndicator(opt.Alias)

	opt.Name = strings.Trim(opt.Name, "-")
	opt.Alias = strings.Trim(opt.Alias, "-")
	opt = updateAttributes(opt, description, toolName)

	return opt
}

func getIndicator(name string) string {
	if strings.Contains(name, "--") {
		return "--"
	}

	if strings.Contains(name, "-") {
		return "-"
	}

	return "-"
}