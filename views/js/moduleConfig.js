var g_financeModule = {
//	lang: "zh-Hans",
	lang: "zh-HANT-TW",
//	filter: 'raw',
	modules: {
        "papersns-form": {
            fullpath: '/app/FormJS',
            requires: ['node', 'widget-base', 'widget-htmlparser', 'io-form', 'widget-parent', 'widget-child', 'base-build', 'substitute', 'io-upload-iframe', 'collection']
        },
        "papersns-form-quickedit": {
        	fullpath: '/app/comboview?js/form/p-form-quickedit.js',
            requires: ["datatable-base", "papersns-form"]
        },
        "finance-module": {
        	fullpath: '/app/comboview?js/financeModule.js',
        	requires: ["papersns-form", "papersns-form-quickedit", "node","widget-base","widget-htmlparser","io-form",
        	           "widget-parent","widget-child","base-build","substitute","io-upload-iframe","collection","overlay",
        	           "calendar","datatype-date","event","json","datatable","datasource-get","datasource-jsonschema",
        	           "datatable-datasource","cssbutton",
        	           "gallery-datatable-paginator","listtemplate-paginator","datatype-date-format","io-base","anim",
        	           "array-extras","querystring-stringify","cssfonts","dataschema-json","datasource-io",
        	           "model-sync-rest","gallery-paginator-view","tabview","panel","dd-plugin",
        	           "gallery-layout","gallery-quickedit","gallery-paginator","datatable-base"]
        },
        "index-module": {
        	fullpath: '/app/comboview?js/financeModule.js',
        	requires: ["node","widget-base",
        	           "widget-parent","widget-child","base-build",
        	           "event","panel","dd-plugin",
        	           "gallery-layout"]
        },
        "step-module": {
        	fullpath: '/app/comboview?js/financeModule.js',
        	requires: ["node","event","querystring-stringify","json","io-form","io-base"]
        }
    }
}
/*
"papersns-form", "papersns-form-quickedit", "node","widget-base","widget-htmlparser","io-form",
        	           "widget-parent","widget-child","base-build","substitute","io-upload-iframe","collection","overlay",
        	           "calendar","datatype-date","event","json","datatable","datasource-get","datasource-jsonschema",
        	           "datatable-datasource","cssbutton",
        	           "gallery-datatable-paginator","listtemplate-paginator","datatype-date-format","io-base","anim",
        	           "array-extras","querystring-stringify","cssfonts","dataschema-json","datasource-io",
        	           "model-sync-rest","gallery-paginator-view","tabview","panel","dd-plugin",
        	           "gallery-layout","gallery-quickedit","gallery-paginator","datatable-base"
 */
var g_financeModuleWithouTable = {
	modules: {
        "papersns-form": {
            fullpath: '/app/FormJS',
            requires: ['node', 'widget-base', 'widget-htmlparser', 'io-form', 'widget-parent', 'widget-child', 'base-build', 'substitute', 'io-upload-iframe', 'collection']
        },
        "papersns-form-quickedit": {
        	fullpath: '/app/comboview?js/form/p-form-quickedit.js',
            requires: ["datatable-base", "papersns-form"]
        },
        "finance-module": {
        	fullpath: '/app/comboview?js/financeModule.js',
        	requires: ["papersns-form", "node","widget-base","widget-htmlparser","io-form",
        	           "widget-parent","widget-child","base-build","substitute","io-upload-iframe","collection","overlay",
        	           "calendar","datatype-date","event","json",
        	           "cssbutton",
        	           "gallery-paginator-view",
        	           
        	           "gallery-quickedit","gallery-paginator",
        	           "gallery-datatable-paginator","datatype-date-format","io-base","anim",
        	           "array-extras","querystring-stringify","cssfonts","dataschema-json","datasource-io",
        	           "model-sync-rest","tabview","panel","dd-plugin"]
        }
    }
}
// withouttable,papersns-form:9.045,
// withouttable,without papersns-form:8.747
// withouttable,only papersns-form and node:8.639
// withouttable,nothing:3.211
//withouttable,papersns-form, but without /app/FormJS:3.558
//withouttable,papersns-form, and all other js, but without /app/FormJS:9.339
//withouttable,papersns-form, node - json, but without /app/FormJS:3.447
//withouttable,papersns-form, node - datasource-io, but without /app/FormJS:3.975
//withouttable,papersns-form, node - dd-plugin, but without /app/FormJS:3.863
//withouttable,papersns-form, node - gallery-layout, but without /app/FormJS:9.096
//withouttable,papersns-form, node - gallery-layout, add gallery-paginator:4.358
//withouttable,papersns-form, node - gallery-layout, add "gallery-paginator-view":4.756
//withouttable,papersns-form, node - gallery-layout, add "gallery-quickedit","gallery-paginator":4.374
//gallery-layout以及,/app/FormJS,都会带来性能的大幅下降,如果我一开始就把所有的JS搞进去?
/*
"datatable-sort","datatable-scroll","resize-plugin",
 */

/*
"papersns-form", "node","widget-base","widget-htmlparser","io-form",
        	           "widget-parent","widget-child","base-build","substitute","io-upload-iframe","collection","overlay",
        	           "calendar","datatype-date","event","json",
        	           "cssbutton",
        	           "datatype-date-format","io-base","anim",
        	           "array-extras","querystring-stringify","cssfonts","dataschema-json","datasource-io",
        	           "model-sync-rest","tabview","panel","dd-plugin",
        	           "gallery-layout" 
 */
