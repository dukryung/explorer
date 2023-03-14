module.exports = {
  parserOptions: {
    ecmaVersion: 2019,
    sourceType: 'module',
  },
  env: {
    es6: true,
    browser: true,
  },
  plugins: [
    'svelte3',
  ],
  extends: [
    'google',
  ],
  settings: {
    'import/resolver': {
      webpack: {
        config: './webpack/common.js',
      },
    },
  },
  overrides: [
    {
      files: ['*.svelte'],
      processor: 'svelte3/svelte3',
    },
  ],
  rules: {
    'require-jsdoc': 0,
    'camelcase': 'off',
    'max-len': [
      'error', {
        code: 120,
      },
    ],
    'valid-jsdoc': [
      'error',
      {
        requireParamDescription: false,
        requireReturnDescription: false,
        requireReturn: false,
        preferType: {
          Boolean: 'boolean',
          Number: 'number',
          object: 'Object',
          String: 'string',
        },
      },
    ],
  },
};
