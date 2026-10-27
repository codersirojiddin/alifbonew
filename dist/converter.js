(function(root){
  'use strict';
  const apostrophes = "'`´ʻʼ‘’";
  const letter = c => !!c && /\p{L}/u.test(c);
  function convert(input, direction='new') {
    let changes=0;
    // Whitespace-delimited URLs and email addresses are kept intact.
    const text=String(input).split(/(\s+)/u).map(token=>{
      if (/:\/\/|@|^www\./i.test(token)) return token;
      if(direction==='old') {
        const allCaps=/\p{Lu}/u.test(token)&&! /\p{Ll}/u.test(token);
        return token.replace(/[öÖğĞşŞçÇ]/gu,c=>{
          changes++;
          return ({ö:"o‘",Ö:"O‘",ğ:"g‘",Ğ:"G‘",ş:'sh',Ş:allCaps?'SH':'Sh',ç:'ch',Ç:allCaps?'CH':'Ch'})[c];
        });
      }
      let out='';
      for(let i=0;i<token.length;i++) {
        const a=token[i],b=token[i+1]||'',low=a.toLowerCase();
        const upper=a!==low;
        let replacement='';
        if ((low==='o'||low==='g')&&apostrophes.includes(b)&&b) {
          if(letter(token[i+2])||(low==='g'&&letter(token[i-1]))) replacement=low==='o'?(upper?'Ö':'ö'):(upper?'Ğ':'ğ');
        } else if(b.toLowerCase()==='h'&&(low==='s'||low==='c')) replacement=low==='s'?(upper?'Ş':'ş'):(upper?'Ç':'ç');
        if(replacement){out+=replacement; changes++; i++;}
        else if(apostrophes.includes(a)&&a!=='ʼ'&&letter(token[i-1])&&letter(b)){out+='ʼ';changes++;}
        else out+=a;
      }
      return out;
    }).join('');
    return {text,changes};
  }
  if(typeof module==='object')module.exports={convert}; else root.Harf={convert};
})(typeof globalThis==='object'?globalThis:this);
