module.exports = {
  modulePathIgnorePatterns: ["<rootDir>/dist/"],
  transform: {
    "^.+\\.(ts|tsx)?$": "ts-jest",
    "^.+\\.(js|jsx)$": "babel-jest",
  },
};
