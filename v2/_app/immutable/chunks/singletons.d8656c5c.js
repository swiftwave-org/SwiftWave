import{w as u}from"./index.4bc12d6c.js";import"./paths.132c185b.js";const A="sveltekit:snapshot",R="sveltekit:scroll",y="sveltekit:index",c={tap:1,hover:2,viewport:3,eager:4,off:-1};function I(e){let t=e.baseURI;if(!t){const o=e.getElementsByTagName("base");t=o.length?o[0].href:e.URL}return t}function x(){return{x:pageXOffset,y:pageYOffset}}const d=new WeakSet,p={"preload-code":["","off","tap","hover","viewport","eager"],"preload-data":["","off","tap","hover"],keepfocus:["","true","off","false"],noscroll:["","true","off","false"],reload:["","true","off","false"],replacestate:["","true","off","false"]};function f(e,t){const o=e.getAttribute(`data-sveltekit-${t}`);return v(e,t,o),o}function v(e,t,o){o!==null&&!d.has(e)&&!p[t].includes(o)&&(console.error(`Unexpected value for ${t} — should be one of ${p[t].map(r=>JSON.stringify(r)).join(", ")}`,e),d.add(e))}const _={...c,"":c.hover};function h(e){let t=e.assignedSlot??e.parentNode;return(t==null?void 0:t.nodeType)===11&&(t=t.host),t}function O(e,t){for(;e&&e!==t;){if(e.nodeName.toUpperCase()==="A"&&e.hasAttribute("href"))return e;e=h(e)}}function U(e,t){let o;try{o=new URL(e instanceof SVGAElement?e.href.baseVal:e.href,document.baseURI)}catch{}const r=e instanceof SVGAElement?e.target.baseVal:e.target,s=!o||!!r||w(o,t)||(e.getAttribute("rel")||"").split(/\s+/).includes("external"),l=(o==null?void 0:o.origin)===location.origin&&e.hasAttribute("download");return{url:o,external:s,target:r,download:l}}function N(e){let t=null,o=null,r=null,s=null,l=null,a=null,n=e;for(;n&&n!==document.documentElement;)r===null&&(r=f(n,"preload-code")),s===null&&(s=f(n,"preload-data")),t===null&&(t=f(n,"keepfocus")),o===null&&(o=f(n,"noscroll")),l===null&&(l=f(n,"reload")),a===null&&(a=f(n,"replacestate")),n=h(n);function i(b){switch(b){case"":case"true":return!0;case"off":case"false":return!1;default:return null}}return{preload_code:_[r??"off"],preload_data:_[s??"off"],keep_focus:i(t),noscroll:i(o),reload:i(l),replace_state:i(a)}}function g(e){const t=u(e);let o=!0;function r(){o=!0,t.update(a=>a)}function s(a){o=!1,t.set(a)}function l(a){let n;return t.subscribe(i=>{(n===void 0||o&&i!==n)&&a(n=i)})}return{notify:r,set:s,subscribe:l}}function k(){const{set:e,subscribe:t}=u(!1);return{subscribe:t,check:async()=>!1}}function w(e,t){return e.origin!==location.origin||!e.pathname.startsWith(t)}function L(e){e.client}const T={url:g({}),page:g({}),navigating:u(null),updated:k()};export{y as I,c as P,R as S,A as a,U as b,N as c,T as d,L as e,O as f,I as g,w as i,x as s};
