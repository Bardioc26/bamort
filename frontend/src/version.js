// Frontend version information
export const VERSION = '0.2.2'

// Git commit will be injected at build time or detected from env
export const GIT_COMMIT = import.meta.env.VITE_GIT_COMMIT || 'unknown'

export function getVersion() {
  return VERSION
}

export function getGitCommit() {
  return GIT_COMMIT
}

export function getVersionInfo() {
  return {
    version: VERSION,
    gitCommit: GIT_COMMIT
  }
}
