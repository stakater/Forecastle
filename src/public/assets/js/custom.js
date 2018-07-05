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
    
    var fileName = "apps.json"
    if(FileExists(fileName)) {
        $.getJSON(fileName, populateAppsListFromJson);
        initSearch();
    } else {
        $.ajax({
            url: "/api/apps",
            type: 'GET',
            success: function(res) {
                $.getJSON(fileName, populateAppsListFromJson);
                initSearch();
            }
        });
    }

    function initSearch() {
        $('.apps .app').each(function(){
            $(this).attr('data-search-term', $(this).find('button span').text().toLowerCase());
        });
    }

    function populateAppsListFromJson(data) {
        $.each(data, renderApp);
    }
    
    function renderApp(key, app) {
        var appTemplate = $("#app-template").html();
        var html = Mustache.render(appTemplate, app)
        $(".apps").append(html);
    }
    
    function FileExists(file)
    {
        var http = new XMLHttpRequest();
        http.open('HEAD', file, false);
        http.send();
        return http.status!=404;
    }
});