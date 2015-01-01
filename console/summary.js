function doRefretorComponent() {
	doRefretor("Component");
}

function doRefretorSelector() {
	doRefretor("Selector");
}

function doRefretorForm() {
	doRefretor("Form");
}

function doRefretorDataSource() {
	doRefretor("DataSource");
}

function doRefretor(name) {
	var dtManager = g_gridPanelDict[name];
	var uri = "/console/refretor?type=" + name;
	YUI(g_financeModule).use("finance-module", function(Y){
		Y.on('io:complete', function(id, o, args) {
			var id = id; // Transaction ID.
			var data = Y.JSON.parse(o.responseText);
			dtManager.dt.set("data", data.items);
			dtManager.hideLoadingImg();
		}, Y, []);
		dtManager.showLoadingImg();
		var request = Y.io(uri);
	});
}
