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
    
        $('.apps .app').each(function(){
            
            if ($(this).filter('[data-search-term *= ' + searchTerm + ']').length > 0 || searchTerm.length < 1) {
                $(this).show();
            } else {
                $(this).hide();
            }
    
        });
    
    });

    var namespacesFilePath = "/etc/cp-config/namespaces.conf";

    $.ajax({
        url: "/file?path=" + namespacesFilePath,
        type: "GET",
        success: function(namespaces) {
            $.ajax({
                url: "/apps?namespaces=" + namespaces,
                type: "GET",
                success: populateAppsListFromJson
            });
        }
    });

    function initSearch() {
        $('.apps .app').each(function(){
            $(this).attr('data-search-term', $(this).find('button>span').text().toLowerCase());
        });
    }

    function populateAppsListFromJson(data) {
        $.each(data, renderApp);
        initSearch();
    }
    
    function renderApp(key, app) {
        var appTemplate = $("#app-template").html();
        var html = Mustache.render(appTemplate, app)
        $(".apps").append(html);
    }
});