var te=Object.defineProperty,ae=Object.defineProperties;var le=Object.getOwnPropertyDescriptors;var E=Object.getOwnPropertySymbols;var ne=Object.prototype.hasOwnProperty,oe=Object.prototype.propertyIsEnumerable;var M=(e,t,a)=>t in e?te(e,t,{enumerable:!0,configurable:!0,writable:!0,value:a}):e[t]=a,S=(e,t)=>{for(var a in t||(t={}))ne.call(t,a)&&M(e,a,t[a]);if(E)for(var a of E(t))oe.call(t,a)&&M(e,a,t[a]);return e},k=(e,t)=>ae(e,le(t));import{Q as se,a as ie,f as re,b as ue,c as D,d as ce,e as R}from"./QInnerLoading.8856fb25.js";import{a as w,d as c,h as g,n as N,F as B,f as x,M as q,$ as de,Q as P,R as I,K as r,I as b,V as A,J as h,W as fe,O as K,H as L,S as V,aJ as z,au as Q,N as G,i as H,s as be,bd as J,T as C,j as F,ac as ge,ae as ve,ai as me,k as pe,e as X,l as ye,r as _e,be as U,aG as he}from"./index.86963db4.js";import"./axios.3f5f2f83.js";import{Q as W}from"./FileSaver.min.d999554e.js";import{g as xe,a as Se,u as $e,d as ke,b as Ce,c as qe,e as Ve}from"./admin.57bb3d7c.js";import"./selection.61faa68a.js";import"./i18n.574036e5.js";const Qe=["top","middle","bottom"];var we=w({name:"QBadge",props:{color:String,textColor:String,floating:Boolean,transparent:Boolean,multiLine:Boolean,outline:Boolean,rounded:Boolean,label:[Number,String],align:{type:String,validator:e=>Qe.includes(e)}},setup(e,{slots:t}){const a=c(()=>e.align!==void 0?{verticalAlign:e.align}:null),l=c(()=>{const i=e.outline===!0&&e.color||e.textColor;return`q-badge flex inline items-center no-wrap q-badge--${e.multiLine===!0?"multi":"single"}-line`+(e.outline===!0?" q-badge--outline":e.color!==void 0?` bg-${e.color}`:"")+(i!==void 0?` text-${i}`:"")+(e.floating===!0?" q-badge--floating":"")+(e.rounded===!0?" q-badge--rounded":"")+(e.transparent===!0?" q-badge--transparent":"")});return()=>g("div",{class:l.value,style:a.value,role:"alert","aria-label":e.label},N(t.default,e.label!==void 0?[e.label]:[]))}});const Be={class:"text-h4"},Ae=B({__name:"SystemInfo",setup(e){const t=x(null);return(async()=>{const l=await xe();t.value=l.result})(),(l,i)=>{var u;return h(),q(A,null,[de("div",Be,[P(I((u=t.value)==null?void 0:u.appName)+" ",1),r(we,{outline:"",align:"middle",color:"primary"},{default:b(()=>{var n;return[P("v"+I((n=t.value)==null?void 0:n.appVersion),1)]}),_:1})]),r(se,{showing:!t.value},{default:b(()=>[r(ie,{size:"xl",color:"primary"})]),_:1},8,["showing"])],64)}}}),Pe=B({__name:"EditSettingsForm",setup(e){const t=fe(),a=x([]),l=x([]);(async()=>{const n=await Se();l.value=n.result,a.value=new Array(l.value.length).fill(!1)})();const u=async(n,f)=>{a.value[n]=!0;try{const o=await $e(l.value[n].key,f);G.create({message:o.message,position:"top",type:"positive"})}finally{a.value[n]=!1}};return(n,f)=>(h(!0),q(A,null,K(l.value,(o,v)=>(h(),q("div",{key:v},[o.isText?(h(),L(W,{key:0,onSubmit:d=>u(v,o.textValue),class:"q-gutter-md flex items-center justify-between"},{default:b(()=>[r(z,{filled:"",type:o.key==="ADMIN_PASSWORD"?"password":"text",modelValue:o.textValue,"onUpdate:modelValue":d=>o.textValue=d,label:`${o.label[V(t).lang.isoName]} *`,class:"tw-flex-1","lazy-rules":"",rules:[d=>d.length>=o.min&&d.length<=o.max||n.$t("forms.rules.rangeLength",{label:`${o.label[V(t).lang.isoName]}`,min:o.min,max:o.max})]},null,8,["type","modelValue","onUpdate:modelValue","label","rules"]),r(Q,{type:"submit",round:"",flat:"",icon:"save",loading:a.value[v]},null,8,["loading"])]),_:2},1032,["onSubmit"])):(h(),L(W,{key:1,onSubmit:d=>u(v,o.numberValue),class:"q-gutter-md flex items-center justify-between"},{default:b(()=>[r(z,{filled:"",type:"number",modelValue:o.numberValue,"onUpdate:modelValue":d=>o.numberValue=d,modelModifiers:{number:!0},label:`${o.label[V(t).lang.isoName]} *`,class:"tw-flex-1","lazy-rules":"",rules:[d=>d>=o.min&&d<=o.max||n.$t("forms.rules.range",{label:`${o.label[V(t).lang.isoName]}`,min:o.min,max:o.max})]},null,8,["modelValue","onUpdate:modelValue","label","rules"]),r(Q,{type:"submit",round:"",flat:"",icon:"save",loading:a.value[v]},null,8,["loading"])]),_:2},1032,["onSubmit"]))]))),128))}}),Ie=["top","right","bottom","left"],Y={type:{type:String,default:"a"},outline:Boolean,push:Boolean,flat:Boolean,unelevated:Boolean,color:String,textColor:String,glossy:Boolean,square:Boolean,padding:String,label:{type:[String,Number],default:""},labelPosition:{type:String,default:"right",validator:e=>Ie.includes(e)},externalLabel:Boolean,hideLabel:{type:Boolean},labelClass:[Array,String,Object],labelStyle:[Array,String,Object],disable:Boolean,tabindex:[Number,String]};function Z(e,t){return{formClass:c(()=>`q-fab--form-${e.square===!0?"square":"rounded"}`),stacked:c(()=>e.externalLabel===!1&&["top","bottom"].includes(e.labelPosition)),labelProps:c(()=>{if(e.externalLabel===!0){const a=e.hideLabel===null?t.value===!1:e.hideLabel;return{action:"push",data:{class:[e.labelClass,`q-fab__label q-tooltip--style q-fab__label--external q-fab__label--external-${e.labelPosition}`+(a===!0?" q-fab__label--external-hidden":"")],style:e.labelStyle}}}return{action:["left","top"].includes(e.labelPosition)?"unshift":"push",data:{class:[e.labelClass,`q-fab__label q-fab__label--internal q-fab__label--internal-${e.labelPosition}`+(e.hideLabel===!0?" q-fab__label--internal-hidden":"")],style:e.labelStyle}}})}}const ee={start:"self-end",center:"self-center",end:"self-start"},Le=Object.keys(ee);var Ne=w({name:"QFabAction",props:k(S({},Y),{icon:{type:String,default:""},anchor:{type:String,validator:e=>Le.includes(e)},to:[String,Object],replace:Boolean}),emits:["click"],setup(e,{slots:t,emit:a}){const l=H(J,()=>({showing:{value:!0},onChildClick:be})),{formClass:i,labelProps:u}=Z(e,l.showing),n=c(()=>{const s=ee[e.anchor];return i.value+(s!==void 0?` ${s}`:"")}),f=c(()=>e.disable===!0||l.showing.value!==!0);function o(s){l.onChildClick(s),a("click",s)}function v(){const s=[];return t.icon!==void 0?s.push(t.icon()):e.icon!==""&&s.push(g(C,{name:e.icon})),(e.label!==""||t.label!==void 0)&&s[u.value.action](g("div",u.value.data,t.label!==void 0?t.label():[e.label])),N(t.default,s)}const d=F();return Object.assign(d.proxy,{click:o}),()=>g(Q,k(S({class:n.value},e),{noWrap:!0,stack:e.stacked,icon:void 0,label:void 0,noCaps:!0,fabMini:!0,disable:f.value,onClick:o}),v)}});const Fe=["up","right","down","left"],Oe=["left","center","right"];var Te=w({name:"QFab",props:k(S(S({},Y),ge),{icon:String,activeIcon:String,hideIcon:Boolean,hideLabel:{default:null},direction:{type:String,default:"right",validator:e=>Fe.includes(e)},persistent:Boolean,verticalActionsAlign:{type:String,default:"center",validator:e=>Oe.includes(e)}}),emits:ve,setup(e,{slots:t}){const a=x(null),l=x(e.modelValue===!0),{proxy:{$q:i}}=F(),{formClass:u,labelProps:n}=Z(e,l),f=c(()=>e.persistent!==!0),{hide:o,toggle:v}=me({showing:l,hideOnRouteChange:f}),d=c(()=>({opened:l.value})),s=c(()=>`q-fab z-fab row inline justify-center q-fab--align-${e.verticalActionsAlign} ${u.value}`+(l.value===!0?" q-fab--opened":" q-fab--closed")),m=c(()=>`q-fab__actions flex no-wrap inline q-fab__actions--${e.direction} q-fab__actions--${l.value===!0?"opened":"closed"}`),p=c(()=>`q-fab__icon-holder  q-fab__icon-holder--${l.value===!0?"opened":"closed"}`);function $(_,O){const T=t[_],j=`q-fab__${_} absolute-full`;return T===void 0?g(C,{class:j,name:e[O]||i.iconSet.fab[O]}):g("div",{class:j},T(d.value))}function y(){const _=[];return e.hideIcon!==!0&&_.push(g("div",{class:p.value},[$("icon","icon"),$("active-icon","activeIcon")])),(e.label!==""||t.label!==void 0)&&_[n.value.action](g("div",n.value.data,t.label!==void 0?t.label(d.value):[e.label])),N(t.tooltip,_)}return pe(J,{showing:l,onChildClick(_){o(_),a.value!==null&&a.value.$el.focus()}}),()=>g("div",{class:s.value},[g(Q,k(S({ref:a,class:u.value},e),{noWrap:!0,stack:e.stacked,align:void 0,icon:void 0,label:void 0,noCaps:!0,fab:!0,"aria-expanded":l.value===!0?"true":"false","aria-haspopup":"true",onClick:v}),y),g("div",{class:m.value},X(t.default))])}});const je={position:{type:String,default:"bottom-right",validator:e=>["top-right","top-left","bottom-right","bottom-left","top","right","bottom","left"].includes(e)},offset:{type:Array,validator:e=>e.length===2},expand:Boolean};function Ee(){const{props:e,proxy:{$q:t}}=F(),a=H(ye,()=>{console.error("QPageSticky needs to be child of QLayout")}),l=c(()=>{const s=e.position;return{top:s.indexOf("top")>-1,right:s.indexOf("right")>-1,bottom:s.indexOf("bottom")>-1,left:s.indexOf("left")>-1,vertical:s==="top"||s==="bottom",horizontal:s==="left"||s==="right"}}),i=c(()=>a.header.offset),u=c(()=>a.right.offset),n=c(()=>a.footer.offset),f=c(()=>a.left.offset),o=c(()=>{let s=0,m=0;const p=l.value,$=t.lang.rtl===!0?-1:1;p.top===!0&&i.value!==0?m=`${i.value}px`:p.bottom===!0&&n.value!==0&&(m=`${-n.value}px`),p.left===!0&&f.value!==0?s=`${$*f.value}px`:p.right===!0&&u.value!==0&&(s=`${-$*u.value}px`);const y={transform:`translate(${s}, ${m})`};return e.offset&&(y.margin=`${e.offset[1]}px ${e.offset[0]}px`),p.vertical===!0?(f.value!==0&&(y[t.lang.rtl===!0?"right":"left"]=`${f.value}px`),u.value!==0&&(y[t.lang.rtl===!0?"left":"right"]=`${u.value}px`)):p.horizontal===!0&&(i.value!==0&&(y.top=`${i.value}px`),n.value!==0&&(y.bottom=`${n.value}px`)),y}),v=c(()=>`q-page-sticky row flex-center fixed-${e.position} q-page-sticky--${e.expand===!0?"expand":"shrink"}`);function d(s){const m=X(s.default);return g("div",{class:v.value,style:o.value},e.expand===!0?m:[g("div",m)])}return{$layout:a,getStickyContent:d}}var Me=w({name:"QPageSticky",props:je,setup(e,{slots:t}){const{getStickyContent:a}=Ee();return()=>a(t)}});var De=(e,t)=>{const a=e.__vccOpts||e;for(const[l,i]of t)a[l]=i;return a};const Re=B({__name:"AdminActions",setup(e){const t=[{icon:"delete_sweep",label:"buttons.deleteExpiredChunks",handler:ke},{icon:"delete",label:"buttons.deleteExpiredFileRecords",handler:Ce},{icon:"delete_outline",label:"buttons.deleteExpiredTextRecords",handler:qe},{icon:"logout",label:"buttons.logout",handler:Ve,callback:()=>{_e.push({name:"AdminLogin"})}}],a=x(new Array(t.length).fill(!1)),l=async(i,u,n)=>{a.value[i]=!0;try{const f=await u();G.create({message:f.message,position:"top",type:"positive"}),n&&n()}finally{a.value[i]=!1}};return(i,u)=>(h(),L(Me,{position:"bottom-left",offset:[18,18]},{default:b(()=>[r(Te,{color:"primary",direction:"up"},{icon:b(({opened:n})=>[r(C,{class:U({"example-fab-animate--hover":n!==!0}),name:"keyboard_arrow_up"},null,8,["class"])]),"active-icon":b(({opened:n})=>[r(C,{class:U({"example-fab-animate":n===!0}),name:"close"},null,8,["class"])]),default:b(()=>[(h(),q(A,null,K(t,(n,f)=>r(Ne,{key:f,onClick:o=>l(f,n.handler,n.callback),loading:a.value[f],color:"primary","external-label":""},{icon:b(()=>[r(C,{name:n.icon},null,8,["name"])]),label:b(()=>[P(I(i.$t(n.label)),1)]),_:2},1032,["onClick","loading"])),64))]),_:1})]),_:1}))}});var ze=De(Re,[["__scopeId","data-v-5be6e34a"]]);const Ze=B({__name:"ConsolePage",setup(e){const t=x("systemInfo");return(a,l)=>(h(),q(A,null,[r(he,null,{default:b(()=>[r(re,{"model-value":20,disable:""},{before:b(()=>[r(ue,{modelValue:t.value,"onUpdate:modelValue":l[0]||(l[0]=i=>t.value=i),vertical:"",class:"text-primary"},{default:b(()=>[r(D,{name:"systemInfo",icon:"info",label:a.$t("tabs.systemInfo")},null,8,["label"]),r(D,{name:"settings",icon:"settings",label:a.$t("tabs.settings")},null,8,["label"])]),_:1},8,["modelValue"])]),after:b(()=>[r(ce,{modelValue:t.value,"onUpdate:modelValue":l[1]||(l[1]=i=>t.value=i),animated:"",vertical:"","keep-alive":"","transition-prev":"jump-up","transition-next":"jump-up"},{default:b(()=>[r(R,{name:"systemInfo"},{default:b(()=>[r(Ae)]),_:1}),r(R,{name:"settings"},{default:b(()=>[r(Pe)]),_:1})]),_:1},8,["modelValue"])]),_:1})]),_:1}),r(ze)],64))}});export{Ze as default};