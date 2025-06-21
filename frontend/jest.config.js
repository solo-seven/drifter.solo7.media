const nextJest = require('next/jest')
const createJestConfig = nextJest({ dir: './' })

const customJestConfig = {
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  testEnvironment: 'jest-environment-jsdom',
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageThreshold: {
    global: { lines: 80, branches: 80, functions: 80, statements: 80 }
  }
}

module.exports = createJestConfig(customJestConfig)
