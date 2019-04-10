export function  formatTime (timestamp) {
    var time0 = new Date(timestamp);
    return time0.toLocaleDateString().replace(/\//g, "-") + " " + time0.toTimeString().substr(0, 8)
}