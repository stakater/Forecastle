/* ------------------------------------------------------------------------------
 *
 *  # Custom JS code
 *
 *  Place here all your custom js. Make sure it's loaded after app.js
 *
 * ---------------------------------------------------------------------------- */

jQuery(document).ready(function($){
    
    $('.search-box').on('keyup', function(){
    
    var searchTerm = $(this).val().toLowerCase();
    
        $('.apps .app').each(function() {
            
            if ($(this).filter('[data-search-term *= ' + searchTerm + ']').length > 0 || searchTerm.length < 1) {
                $(this).show();
            } else {
                $(this).hide();
            }
    
        });

        $('.group').each(function() {
            if ($(this).find(".app").filter(":visible").length == 0) {
                $(this).find(".group-header").hide();
            } else {
                $(this).find(".group-header").show();
            }
        });
    
    });

    $.ajax({
        url: "/apps",
        type: "GET",
        success: populateAppsListFromJson,
        complete: getConfigProperties
    });

    function initSearch() {
        $('.apps .app').each(function(){
            $(this).attr('data-search-term', $(this).find('button>span').text().toLowerCase());
        });
    }

    function populateAppsListFromJson(data) {
        data = groupByKeys(data, "Group", "Apps");
        $.each(data, renderApp);
        initSearch();
    }

    function getConfigProperties() {
        $.ajax({
            url: "/config",
            type: "GET",
            success: renderConfigProperties,
            complete: function() {
                $(".loader").fadeOut(500);
            }
        });
    }

    function renderConfigProperties(config) {
        if(config.HeaderBackground !== "") {
            $(".page-header").css("background-color", config.HeaderBackground);
        }
        if(config.HeaderForeground !== "") {
            $(".page-header").css("color", config.HeaderForeground);
        }
        if(config.Title !== "") {
            $(".page-title h2").html(config.Title);
        }
    }
    
    function renderApp(key, app) {
        var appTemplate = $("#app-template").html();
        var template = Handlebars.compile(appTemplate)
        var html = template(app)
        $(".apps").append(html);
    }

    var groupByKeys = function(data, key, valuesKey) {
        var map = new Map();

        data.forEach(item => {
            if(!map.has(item[key])) {
                map.set(item[key], [])
            }
            map.get(item[key]).push(item)
        });

        var result = [];

        map.forEach((mValue, mKey) => {
            result.push({
                [key]: mKey,
                [valuesKey]: mValue
            })
        })

        return result
    };
});