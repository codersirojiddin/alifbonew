const assert=require('node:assert/strict');
const {convert}=require('../dist/converter.js');
const cases=[
 ["O‘zbekiston, g‘alaba, shahar, chiroyli.",'Özbekiston, ğalaba, şahar, çiroyli.'],
 ["O'ZBEK SHAHAR CHIROQ TOG'",'ÖZBEK ŞAHAR ÇIROQ TOĞ'],
 ["o`rik gʻisht oʼq g’oz o´rmon",'örik ğişt öq ğoz örmon'],
 ["ma'no san’at Isʼhoq",'maʼno sanʼat Isʼhoq'],
 ['https://example.com/sh www.shahar.uz user@shahar.uz shahar','https://example.com/sh www.shahar.uz user@shahar.uz şahar'],
 ["'o' va 'g'", "'o' va 'g'"],['',''],['ÖĞŞÇ ö ğ ş ç','ÖĞŞÇ ö ğ ş ç'],
 ['Shahar — SHAHAR — Chiroq — CHIROQ','Şahar — ŞAHAR — Çiroq — ÇIROQ']
];
for(const [input,expected]of cases) assert.equal(convert(input).text,expected,input);
assert.equal(convert('Özbekiston, Şahar, ŞAHAR, ÇIROQ','old').text,'O‘zbekiston, Shahar, SHAHAR, CHIROQ');
assert.equal(convert('ÖĞŞÇ','old').text,'O‘G‘SHCH');
assert.equal(convert("o‘g‘shch").changes,4);
assert.equal(convert('Oddiy matn').changes,0);
const long="O‘zbekiston shahar chiroq tog‘ ".repeat(10000);
assert.equal(convert(long).text,"Özbekiston şahar çiroq toğ ".repeat(10000));
console.log('15 conversion checks passed, including 300k-character input.');
