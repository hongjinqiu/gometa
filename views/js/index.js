var TAB_ID = 100;
var TAB_DICT = {};// id:title
var TAB_DICT_REVERSE = {};// title:id,

/**
 * 新增tab,清除 tabs-selected,为新增上的tab添加 tabs-selected
 * @param title
 * @param url
 */
function putTabIfAbsent(title, url, isRefresh) {
	if (TAB_DICT_REVERSE[title]) {
		if (isRefresh) {
			closeTab(TAB_DICT_REVERSE[title]);
			putTabIfAbsent(title, url);
		} else {
			selectTab(TAB_DICT_REVERSE[title]);
		}
	} else {
		var id = TAB_ID++;
		var tabId = "tab" + id;
		var tabsWrapper = document.getElementById("tabsWrapper");
		var liElem = document.createElement("li");
		liElem.className = "tabs-selected";
		var liHtmlLi = [];
		liHtmlLi.push("<a class=\"tabs-inner\" href=\"javascript:void(0)\" onclick=\"selectTab('{tabId}')\">");
		liHtmlLi.push("<span class=\"tabs-title\">{title}</span>");
		liHtmlLi.push("<span class=\"tabs-icon\"></span>");
		liHtmlLi.push("</a>");
		liHtmlLi.push("<a class=\"tabs-close\" onclick=\"closeTab('{tabId}')\" href=\"javascript:void(0)\"></a>");
		var liHtml = liHtmlLi.join("");
		liHtml = liHtml.replace(/{tabId}/gi, tabId);
		liHtml = liHtml.replace(/{title}/gi, title);
		liElem.id = tabId;
		liElem.innerHTML = liHtml;
		
		// 子iframe
		var tabsChildWrapper = document.getElementById("tabsChildWrapper");
		var childElem = document.createElement("div");
		var childHtmlLi = [];
		childHtmlLi.push("	<iframe width=\"100%\" height=\"100%\" frameborder=\"0\" src=\"{url}\"></iframe>");
		var childHtml = childHtmlLi.join("");
		childHtml = childHtml.replace(/{url}/gi, url);
		childElem.id = tabId + "Child";
		childElem.className = "tabChild";
		childElem.innerHTML = childHtml;

		hideCurrentSelectTab();
		
		TAB_DICT_REVERSE[title] = tabId;
		TAB_DICT[tabId] = title;
		tabsWrapper.appendChild(liElem);
		tabsChildWrapper.appendChild(childElem);
	}
}

function hideCurrentSelectTab() {
	var tabsWrapper = document.getElementById("tabsWrapper");
	var tabsChildLi = tabsWrapper.getElementsByTagName("li");
	for (var i = 0; i < tabsChildLi.length; i++) {
		if (tabsChildLi[i].className.indexOf("tabs-selected") > -1) {
			tabsChildLi[i].className = tabsChildLi[i].className.replace("tabs-selected", "");
			document.getElementById(tabsChildLi[i].id + "Child").style.display = "none";
			break;
		}
	}
}

function getCurrentSelectTab() {
	var tabsWrapper = document.getElementById("tabsWrapper");
	var tabsChildLi = tabsWrapper.getElementsByTagName("li");
	for (var i = 0; i < tabsChildLi.length; i++) {
		if (tabsChildLi[i].className.indexOf("tabs-selected") > -1) {
			return tabsChildLi[i];
		}
	}
	return null;
}

function closeTab(tabId) {
	// 关闭,同时选中前一个,
	var tabsWrapper = document.getElementById("tabsWrapper");
	var tabsChildWrapper = document.getElementById("tabsChildWrapper");
	
	var tabToRemove = document.getElementById(tabId);
	var tabChildToRemove = document.getElementById(tabId + "Child");
	
	var nextTabId = getNextTabId(tabId);
	
	tabsWrapper.removeChild(tabToRemove);
	tabsChildWrapper.removeChild(tabChildToRemove);
	var title = TAB_DICT[tabId];
	delete TAB_DICT[tabId];
	delete TAB_DICT_REVERSE[title];
	
	if (nextTabId) {
		selectTab(nextTabId);
	} else {// 选中最后一个
		var li = tabsWrapper.getElementsByTagName("li");
		selectTab(li[li.length - 1].id);
	}
}

function getNextTabId(tabId) {
	var tabsWrapper = document.getElementById("tabsWrapper");
	var tabsChildLi = tabsWrapper.getElementsByTagName("li");
	for (var i = 0; i < tabsChildLi.length; i++) {
		if (tabsChildLi[i].id == tabId) {
			if (i != (tabsChildLi.length - 1)) {
				return tabsChildLi[i+1].id;
			}
		}
	}
	return null;
}

function selectTab(tabId) {
	var currentSelectTab = getCurrentSelectTab();
	if (!currentSelectTab || (currentSelectTab && currentSelectTab.id != tabId)) {
		hideCurrentSelectTab();
		
		document.getElementById(tabId).className += " tabs-selected";
		document.getElementById(tabId + "Child").style.display = "";
	}
}

