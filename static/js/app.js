let all_page = parseInt("{{.AllPage}}");
let page = parseInt("{{.Page}}");
let count = parseInt("{{.Count}}");
const pre = document.getElementById('pre');
const next = document.getElementById('next');
if (Math.ceil(count / 5) > page) {
    next.innerText = "下一页";
    next.href = "/?page=" + (page + 1)
}
if (page !== 1) {
    pre.innerText = "上一页"
    pre.href = "/?page=" + (page - 1)
}