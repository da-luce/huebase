package vscode

import (
	"encoding/json"

	"github.com/da-luce/paletteport/internal/color"
)

type VSCodeTheme struct {
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	Maintainers   []string `json:"maintainers"`
	Type          string   `json:"type"`
	SemanticClass string   `json:"semanticClass"`
	Colors        Colors   `json:"colors"`
}

type Color = color.Color
type Colors struct {
	// Top‑level flat keys
	Foreground            *Color `json:"foreground"`
	DescriptionForeground *Color `json:"descriptionForeground"`
	DisabledForeground    *Color `json:"disabledForeground"`
	FocusBorder           *Color `json:"focusBorder"`
	ErrorForeground       *Color `json:"errorForeground"`
	WidgetShadow          *Color `json:"widget.shadow"`
	ScrollbarShadow       *Color `json:"scrollbar.shadow"`
	BadgeBackground       *Color `json:"badge.background"`
	BadgeForeground       *Color `json:"badge.foreground"`
	IconForeground        *Color `json:"icon.foreground"`
	SettingsHeader        *Color `json:"settings.headerForeground"`
	WindowActiveBorder    *Color `json:"window.activeBorder"`
	WindowInactiveBorder  *Color `json:"window.inactiveBorder"`
	SashHoverBorder       *Color `json:"sash.hoverBorder"`
	ToolbarActiveBG       *Color `json:"toolbar.activeBackground"`
	ToolbarHoverBG        *Color `json:"toolbar.hoverBackground"`

	ExtensionButtonProminentBackground *Color `json:"extensionButton.prominentBackground"`
	ExtensionButtonProminentHover      *Color `json:"extensionButton.prominentHoverBackground"`
	ExtensionButtonProminentForeground *Color `json:"extensionButton.prominentForeground"`
	ExtensionBadgeRemoteBackground     *Color `json:"extensionBadge.remoteBackground"`
	ExtensionBadgeRemoteForeground     *Color `json:"extensionBadge.remoteForeground"`

	ButtonBackground          *Color `json:"button.background"`
	ButtonHoverBackground     *Color `json:"button.hoverBackground"`
	ButtonSecondaryBackground *Color `json:"button.secondaryBackground"`
	ButtonForeground          *Color `json:"button.foreground"`
	ProgressBarBackground     *Color `json:"progressBar.background"`

	InputBackground             *Color `json:"input.background"`
	InputForeground             *Color `json:"input.foreground"`
	InputBorder                 *Color `json:"input.border"`
	InputPlaceholderForeground  *Color `json:"input.placeholderForeground"`
	InputOptionActiveForeground *Color `json:"inputOption.activeForeground"`
	InputOptionActiveBackground *Color `json:"inputOption.activeBackground"`

	InputValidationInfoForeground    *Color `json:"inputValidation.infoForeground"`
	InputValidationInfoBackground    *Color `json:"inputValidation.infoBackground"`
	InputValidationInfoBorder        *Color `json:"inputValidation.infoBorder"`
	InputValidationWarningForeground *Color `json:"inputValidation.warningForeground"`
	InputValidationWarningBackground *Color `json:"inputValidation.warningBackground"`
	InputValidationWarningBorder     *Color `json:"inputValidation.warningBorder"`
	InputValidationErrorForeground   *Color `json:"inputValidation.errorForeground"`
	InputValidationErrorBackground   *Color `json:"inputValidation.errorBackground"`
	InputValidationErrorBorder       *Color `json:"inputValidation.errorBorder"`

	DropdownForeground     *Color `json:"dropdown.foreground"`
	DropdownBackground     *Color `json:"dropdown.background"`
	DropdownListBackground *Color `json:"dropdown.listBackground"`

	TreeIndentGuidesStroke *Color `json:"tree.indentGuidesStroke"`

	// Nested objects
	ActivityBar      ActivityBarColors    `json:"activityBar"`
	ActivityBarBadge BadgeColors          `json:"activityBarBadge"`
	ActivityBarTop   ActivityBarTopColors `json:"activityBarTop"`

	SideBar                SideBarColors       `json:"sideBar"`
	SideBarTitleForeground *Color              `json:"sideBarTitle.foreground"`
	SideBarSectionHeader   SideBarHeaderColors `json:"sideBarSectionHeader"`
	SideBarDropBackground  *Color              `json:"sideBar.dropBackground"`
}

// Example nested structs
type ActivityBarColors struct {
	Background         *Color `json:"background"`
	Foreground         *Color `json:"foreground"`
	InactiveForeground *Color `json:"inactiveForeground"`
	Border             *Color `json:"border"`
}

type BadgeColors struct {
	Background *Color `json:"background"`
	Foreground *Color `json:"foreground"`
}

type ActivityBarTopColors struct {
	Foreground         *Color `json:"foreground"`
	InactiveForeground *Color `json:"inactiveForeground"`
}

type SideBarColors struct {
	Background *Color `json:"background"`
	Foreground *Color `json:"foreground"`
	Border     *Color `json:"border"`
}

type SideBarHeaderColors struct {
	Background *Color `json:"background"`
	Foreground *Color `json:"foreground"`
	Border     *Color `json:"border"`
}

func (rw *VSCodeTheme) FromString(input string) error {
	err := json.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}
