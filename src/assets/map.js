import {LeafletMap, TileLayer, Marker} from 'leaflet';
import {MarkerClusterGroup} from 'leaflet.markercluster';

function initializeMap() {
    const map = new LeafletMap('map').setView([60.2, 20.2], 10);

    new TileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
    }).addTo(map);

    return map;
}

function addMarkers(map) {
    const locations = JSON.parse(document.getElementById('map-container').dataset.locations);
    // Tooltips are hover-only, on touch screens the popup is the way to reach a listing
    const canHover = window.matchMedia('(hover: hover)').matches;

    const cluster = new MarkerClusterGroup();

    locations
        .filter(l => l.Lat != undefined && l.Lon != undefined && l.Lat.Valid && l.Lon.Valid)
        .forEach(l => {
            const marker = new Marker([l.Lat.Float64, l.Lon.Float64]);

            const link = document.createElement('a');
            link.href = encodeURI(l.Url);
            link.textContent = l.Address;
            marker.bindPopup(link);

            if (canHover) {
                const address = document.createElement('span');
                address.textContent = l.Address;
                marker.bindTooltip(address);
            }

            cluster.addLayer(marker);
        });

    map.addLayer(cluster);
}

const map = initializeMap();
addMarkers(map);
