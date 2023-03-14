import {writable} from 'svelte/store';

export const stats = writable({
  block_height: 0,
  tx_total: 0,
  token_total: 0,
  total_bonded_tokens: 0,
  validator_total: 0,
  block_avg_time: 0,
  block_min_time: 0,
});

export const networkConnection = writable(true);

export const localeGlobalState = writable('en-US');