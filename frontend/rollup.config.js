export default {
  // 'external' tells rollup.js not to generate a warning when
  // those libraries are used. It indicates that they are expected
  // to be external dependencies loaded directly from the html.
  external: ["react", "react-dom"],
  output: {
    // 'globals' tells rollup.js that the specified libraries are
    // available as global javascript values before the code
    // is run. Every external dependency loaded directly from the
    // html should be listed as a global.
    globals: {
      react: "React",
      "react-dom": "ReactDOM",
    },
  },
  plugins: [],
};
