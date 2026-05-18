package p_totschool_appointments

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

var (
	pluginPagePatches []registry.Pair[string, func(components.PageInterface) components.PageInterface]
	pluginViewPatches []registry.Pair[string, func(*views.View) *views.View]
)

func patchPluginPage(key string, patch func(components.PageInterface) components.PageInterface) {
	pluginPagePatches = append(pluginPagePatches, registry.Pair[string, func(components.PageInterface) components.PageInterface]{
		Key: key, Value: patch,
	})
}

func patchPluginView(key string, patch func(*views.View) *views.View) {
	pluginViewPatches = append(pluginViewPatches, registry.Pair[string, func(*views.View) *views.View]{
		Key: key, Value: patch,
	})
}

func pluginPagesWithPatches(entries []registry.Pair[string, components.PageInterface]) lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: entries,
		Patches: pluginPagePatches,
	}
}

func pluginViewsWithPatches(entries []registry.Pair[string, *views.View]) lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: entries,
		Patches: pluginViewPatches,
	}
}
