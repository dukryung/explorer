const {GoogleSpreadsheet} = require('google-spreadsheet');
const fs = require('fs');
const cert = require('./google-auth.json');


async function fetch() {
  const id = '1aYS_wpVQIq5jvgID6Apdy32964XefyVfbk_JhTiAI0o';

  const doc = new GoogleSpreadsheet(id);
  await doc.useServiceAccountAuth(cert);
  await doc.loadInfo();

  const i18nSheet = doc.sheetsByIndex[0];
  await i18nSheet.loadCells('A1:D999');

  const keys = await i18nKeys(i18nSheet);

  for (let i = 0; i < keys.length; i++) {
    const i18n = i18nParse(i18nSheet, i+1);
    i18nSave(i18n, keys[i]);
  }

  return new Promise(function(resolve, reject) {});
}

function i18nSave(i18n, key) {
  const json = JSON.stringify(i18n, null, '\t');
  const path = './src/i18n/';
  const file = key + '.json';
  fs.writeFile(path + file, json, function(err) {
    console.log('fetch ' + file);
  });
}

function i18nParse(sheet, idx) {
  let data = {};
  for (let row = 1; row < sheet.rowCount; row++) {
    const key = sheet.getCell(row, 0).value;
    const value = sheet.getCell(row, idx).value;

    if (key == null) {
      break;
    }

    const cur = [];
    cur[key] = value;

    data = Object.assign(data, cur);
  }

  return data;
}

// get i18n keys
function i18nKeys(sheet) {
  const keys = [];
  for (let i = 1; i < sheet.columnCount; i++) {
    keys.push(sheet.getCell(0, i).value);
  }

  return keys;
}

fetch().then();
