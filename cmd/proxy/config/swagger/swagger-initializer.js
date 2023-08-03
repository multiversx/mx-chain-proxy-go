window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // Custom plugin to hide the API definition URL
  const HideInfoUrlPartsPlugin = () => {
    return {
      wrapComponents: {
        InfoUrl: () => () => null
      }
    }
  }

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "openapi.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset.slice(1)
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl,
      HideInfoUrlPartsPlugin
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
