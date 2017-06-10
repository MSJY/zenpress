package hook

import (
	"fmt"
	"pueue"
	//"sync/atomic"
)

var (
	KueuesMap         = map[string]*pueue.PriorityQueue{}
	FiltersMap        = map[string][]byte{}
	ops        uint64 = 0

	CurrentFilterKey string
)

type Node struct {
	priority uint64
	key      string
	callback func([]byte) []byte
}

func (n *Node) Less(other interface{}) bool {
	return n.priority < other.(*Node).priority
}

func RemoveFilter(key string) {
	if HasFilter(key) {
		delete(FiltersMap, key)
	}
}

func RemoveAction(key string) {
	if HasAction(key) {
		delete(KueuesMap, key)
	}
}

func RemoveActions() {
	KueuesMap = map[string]*pueue.PriorityQueue{}
}

func HasFilter(key string) bool {
	if _, okay := FiltersMap[key]; okay {
		return true
	}
	return false
}

func HasAction(key string) bool {
	if _, okay := KueuesMap[key]; okay {
		return true
	}
	return false
}

func AddAction(key string, function func(), priorities ...uint64) {
	AddFilter(key, func([]byte) []byte {
		function()
		return nil
	}, priorities...)
}

func AddFilter(key string, function func([]byte) []byte, priorities ...uint64) {
	if !HasAction(key) {
		KueuesMap[key] = pueue.New()
	}

	var priority uint64
	if len(priorities) > 0 {
		priority = priorities[0]
	} else {
		//atomic.AddUint64(&ops, 1)
		priority = ops
	}

	KueuesMap[key].Push(&Node{priority: priority, key: key, callback: function})
}

func CurrentFilter() string {
	return CurrentFilterKey
}

func DoFilter(key string, function func() []byte) []byte {
	CurrentFilterKey = key
	if HasAction(key) {
		for KueuesMap[key].Len() > 0 {
			n, okay := KueuesMap[key].Pop().(*Node)
			if !okay {
				continue
			}

			fmt.Println("DoFilter:", n.priority, n.key)
			if function != nil {
				if FiltersMap[key] != nil {
					FiltersMap[key] = n.callback(FiltersMap[key])
				} else {
					FiltersMap[key] = n.callback(function())
				}
			} else { //for action
				FiltersMap[key] = n.callback(nil)
			}

		}
	}
	return FiltersMap[key]
}

func DoAction(key string) {
	DoFilter(key, nil)
}

// RegisterActivation Set the activation hook for a plugin.
/*
 * When a plugin is activated, the action 'activate_PLUGINNAME' hook is
 * called. In the name of this hook, PLUGINNAME is replaced with the name
 * of the plugin, including the optional subdirectory. For example, when the
 * plugin is located in wp-content/plugin/sampleplugin/sample.app, then
 * the name of this hook will become 'activate_sampleplugin/sample.app'.
 *
 * When the plugin consists of only one file and is (as by default) located at
 * wp-content/plugin/sample.app the name of this hook will be
 * 'activate_sample.app'.
 *
 * @param string   key    The filename of the plugin including the path.
 * @param callable function The function hooked to the 'activate_PLUGIN' action.
 */
func RegisterActivation(function func(), key ...string) {
	var iKey = "activate"
	if len(key) > 0 {
		iKey = iKey + "_" + key[0]
	}
	AddAction(iKey, function)
}

// RegisterDeactivation Set the deactivation hook for a plugin.
/*
 * When a plugin is deactivated, the action 'deactivate_PLUGINNAME' hook is
 * called. In the name of this hook, PLUGINNAME is replaced with the name
 * of the plugin, including the optional subdirectory. For example, when the
 * plugin is located in content/plugin/sampleplugin/sample.app, then
 * the name of this hook will become 'deactivate_sampleplugin/sample.app'.
 *
 * When the plugin consists of only one file and is (as by default) located at
 * content/plugin/sample.app the name of this hook will be
 * 'deactivate_sample.app'.
 *
 *
 * @param string   key     The filename of the plugin including the path.
 * @param callable function The function hooked to the 'deactivate_PLUGIN' action.
 */
func RegisterDeactivation(function func(), key ...string) {
	var iKey = "deactivate"
	if len(key) > 0 {
		iKey = iKey + "_" + key[0]
	}
	AddAction(iKey, function)
}

//////////////////////////////////

/**
 * Add a submenu page.
 *
 * This function takes a capability which will be used to determine whether
 * or not a page is included in the menu.
 *
 * The function which is hooked in to handle the output of the page must check
 * that the user has the required capability as well.
 *
 * @global arraysubmenu
 * @global arraymenu
 * @global array_wp_real_parent_file
 * @global bool _wp_submenu_nopriv
 * @global array_registered_pages
 * @global array_parent_pages
 *
 * @param string  parent_slug The slug name for the parent menu (or the file name of a standard WordPress admin page).
 * @param string  page_title  The text to be displayed in the title tags of the page when the menu is selected.
 * @param string  menu_title  The text to be used for the menu.
 * @param string  capability  The capability required for this menu to be displayed to the user.
 * @param string  menu_slug   The slug name to refer to this menu by (should be unique for this menu).
 * @param callable function    The function to be called to output the content for this page.
 * @return false|string The resulting page's hook_suffix, or false if the user does not have the capability required.
 */
/*
func  AddSubmenuPage(parent_slug,page_title,menu_title,capability,menu_slug string,function ...func() ) {

	if ( !current_user_can(capability ) ) {
		_wp_submenu_nopriv[parent_slug][menu_slug] = true;
		return false;
	}


	 // If the parent doesn't already have a submenu, add a link to the parent
	 // as the first item in the submenu. If the submenu file is the same as the
	 // parent file someone is trying to link back to the parent manually. In
	 // this case, don't automatically add a link back to avoid duplication.

	if (!isset(submenu[parent_slug] ) &&menu_slug !=parent_slug ) {
		foreach ( (array)menu asparent_menu ) {
			if (parent_menu[2] ==parent_slug && current_user_can(parent_menu[1] ) )
				submenu[parent_slug][] = array_slice(parent_menu, 0, 4 );
		}
	}

	submenu[parent_slug][] = array (menu_title,capability,menu_slug,page_title );

	hookname :=  menu_slug+parent_slug;

if (len(function)>0)&&(len(hookname)>0) {
	AddAction(hookname,function[0])
}
	_registered_pages[hookname] = true;

	// No parent as top level.
	_parent_pages[menu_slug] =parent_slug;

	return hookname;
}
*/

/**
 * Add submenu page to the Settings main menu.
 *
 * This function takes a capability which will be used to determine whether
 * or not a page is included in the menu.
 *
 * The function which is hooked in to handle the output of the page must check
 * that the user has the required capability as well.
 *
 * @param string   page_title The text to be displayed in the title tags of the page when the menu is selected.
 * @param string   menu_title The text to be used for the menu.
 * @param string   capability The capability required for this menu to be displayed to the user.
 * @param string   menu_slug  The slug name to refer to this menu by (should be unique for this menu).
 * @param callable function   The function to be called to output the content for this page.
 * @return false|string The resulting page's hook_suffix, or false if the user does not have the capability required.
 */
/*
func AddOptionsPage(page_title, menu_title, capability, menu_slug string, function ...func()) {
	return AddSubmenuPage("options-general", page_title, menu_title, capability, menu_slug, function...)
}
*/
