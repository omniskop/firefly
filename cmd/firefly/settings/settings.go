package settings

import (
	"reflect"

	"github.com/therecipe/qt/core"
)

var qSettings = core.NewQSettings("omniskop", "firefly", nil)

var changeListeners = make(map[string][]func(interface{}))

// Set the key to value.
// If the new value differs from the old one the registered change listeners will be called.
func Set(key string, value interface{}) {
	if !reflect.DeepEqual(Get(key), value) {
		qSettings.SetValue(key, core.NewQVariant1(value))
		callChangeListeners(key, value)
	}
}

func callChangeListeners(key string, value interface{}) {
	listeners := changeListeners[key]
	for _, listener := range listeners {
		listener(value)
	}
}

// Remove the key and it's change listeners
func Remove(key string) {
	qSettings.Remove(key)
	delete(changeListeners, key)
}

// Get the value of the given key
func Get(key string) interface{} {
	value := qSettings.Value(key, core.NewQVariant())
	if !value.IsValid() {
		return nil
	}
	return value.ToInterface()
}

// GetWithDefault returns the value of the key or the defaultValue when the key is not set
func GetWithDefault(key string, defaultValue interface{}) interface{} {
	return qSettings.Value(key, core.NewQVariant1(defaultValue)).ToInterface()
}

// GetString returns the string that is stored under the key.
// If the key is not set it returns an empty string.
// For more information about value conversion check http://doc.qt.io/qt-5/qvariant.html#toString.
func GetString(key string) string {
	// if the returned QVariant is not of type string the ToString() method returns an empty string
	return qSettings.Value(key, core.NewQVariant15("")).ToString()
}

// GetStrings returns the string slice that is stored under the key.
// If the key is not set it returns an empty slice.
// For more information about value conversion check http://doc.qt.io/qt-5/qvariant.html#toStringList.
func GetStrings(key string) []string {
	// if the returned QVariant is not of type []string and cannot be converted to it the ToStringList() method returns an empty slice
	return qSettings.Value(key, core.NewQVariant17([]string{})).ToStringList()
}

// GetBool returns the bool that is stored under the key.
// If the key is not set it returns false.
// For more information about value conversion check https://doc.qt.io/qt-5/qvariant.html#toBool.
func GetBool(key string) bool {
	return qSettings.Value(key, core.NewQVariant9(false)).ToBool()
}

// GetInt returns the int that is stores under the key.
// If the key is not set it returns 0.
// For more information about value conversion check http://doc.qt.io/qt-5/qvariant.html#toInt.
func GetInt(key string) int {
	return qSettings.Value(key, core.NewQVariant5(0)).ToInt(nil)
}

// OnChange registers the callback to be called when the value of the key changes
func OnChange(key string, callback func(interface{})) {
	listeners := changeListeners[key]
	listeners = append(listeners, callback)
	changeListeners[key] = listeners
}

// RemoveChangeListener removes the given change listener from the key
func RemoveChangeListener(key string, callback func(interface{})) {
	listeners := changeListeners[key]
	callbackPointer := reflect.ValueOf(callback).Pointer()
	for i, listener := range listeners {
		if reflect.ValueOf(listener).Pointer() == callbackPointer {
			c := len(listeners) - 1
			listeners[i] = listeners[c]
			listeners = listeners[:c]
			changeListeners[key] = listeners
			return
		}
	}
}
