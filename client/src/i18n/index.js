import {addMessages, init, getLocaleFromNavigator} from 'svelte-i18n';
import enUS from './en-US.json';
import koKR from './ko-KR.json';
import ruRU from './ru-RU.json';

function addMessage() {
  addMessages('en-US', enUS);
  // addMessages('en', enUS);
  addMessages('ko-KR', koKR);
  // addMessages('ko', koKR);
  // addMessages('ru-RU', ruRU);
  // addMessages('ru', ruRU);
}

export default function i18n() {
  addMessage();

  init({
    fallbackLocale: 'en',
    initialLocale: getLocaleFromNavigator(),
  });
}
