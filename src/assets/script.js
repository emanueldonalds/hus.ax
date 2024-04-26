populateInputFields();
registerUrlUpdater();

function populateInputFields() {
    const query = window.location.search;
    const params = new URLSearchParams(query);

    const agencyName = "agency";
    const priceMinName = "price_min";
    const priceMaxName = "price_max";
    const yearMinName = "year_min";
    const yearMaxName = "year_max";
    const areaMinName = "size_value_min";
    const areaMaxName = "size_value_max";
    const priceOverAreaMinName = "price_over_area_min";
    const priceOverAreaMaxName = "price_over_area_max";
    const includeDeletedName = "include_deleted";
    const orderByName = "order_by";
    const sortOrderName = "sort_order";

    document.getElementById(agencyName).value = params.get(agencyName) ?? '';
    document.getElementById(priceMinName).value = params.get(priceMinName);
    document.getElementById(priceMaxName).value = params.get(priceMaxName);
    document.getElementById(yearMinName).value = params.get(yearMinName);
    document.getElementById(yearMaxName).value = params.get(yearMaxName);
    document.getElementById(areaMinName).value = params.get(areaMinName);
    document.getElementById(areaMaxName).value = params.get(areaMaxName);
    document.getElementById(priceOverAreaMinName).value = params.get(priceOverAreaMinName);
    document.getElementById(priceOverAreaMaxName).value = params.get(priceOverAreaMaxName);
    document.getElementById(includeDeletedName).checked = params.get(includeDeletedName) === "true";
    document.getElementById(orderByName).value = params.get(orderByName);
    document.getElementById(sortOrderName).value = params.get(sortOrderName);
}

function registerUrlUpdater() {
    document.body.addEventListener('htmx:beforeRequest', function(event) {
        const windowSearch = new URLSearchParams(window.location.search);
        const requestParams = Object.entries(event.detail.requestConfig.parameters).map(e => {
            return { key: e[0], value: e[1] }
        });

        const emptyParams = requestParams.filter(p => p.value === "")
        const nonEmptyParams = requestParams.filter(p => p.value !== "")

        for (let param of emptyParams) {
            windowSearch.delete(param.key)
        }
        for (let param of nonEmptyParams) {
            windowSearch.set(param.key, param.value)
        }

        window.history.pushState(null, null, "?" + windowSearch.toString());
    });
}

// Sets sort_by and sort_order and submits the search form
function sort(field) {
    const orderByInput = document.getElementById('order_by')
    const sortOrderInput = document.getElementById('sort_order')

    orderByInput.value = field;

    if (orderByInput.value === field) {
        sortOrderInput.value = sortOrderInput.value === "desc" ? "asc" : "desc";
    }
    else {
        sortOrderInput.value = "desc";
    }
    document.getElementById("send").click();
}
