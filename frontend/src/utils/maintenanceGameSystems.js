import API from './api'

const defaultSystemLabel = (system = {}) => {
  const code = system.code || ''
  const name = system.name || ''
  if (code && name) return `${code} (${name})`
  return code || name || String(system.id ?? '')
}

export const normalizeSystem = (gs = {}) => ({
  ...gs,
  id: gs.id ?? gs.ID ?? gs.Id ?? null,
  code: gs.code ?? gs.Code ?? '',
  name: gs.name ?? gs.Name ?? '',
  description: gs.description ?? gs.Description ?? '',
  is_active: gs.is_active ?? gs.IsActive ?? gs.isActive ?? false,
})

export const buildSystemOptions = (gameSystems = [], labelBuilder = defaultSystemLabel) =>
  gameSystems.map(system => ({
    id: system.id,
    label: labelBuilder(system),
  }))

export const systemOptionsFor = (gameSystems = [], labelBuilder = defaultSystemLabel) =>
  buildSystemOptions(gameSystems, labelBuilder)

export const findSystemById = (gameSystems = [], id) => {
  if (id === null || id === undefined) return null
  return gameSystems.find(gs => gs.id === id) || null
}

export const findSystemIdByCode = (gameSystems = [], code) => {
  if (!code) return null
  const match = gameSystems.find(gs => gs.code === code)
  return match ? match.id : null
}

export const buildGameSystemParams = system => {
  if (!system) return {}
  return {
    game_system_id: system.id,
    game_system: system.name,
  }
}

export const getSystemCodeById = (gameSystems = [], systemId, fallback = '') => {
  if (!systemId) return fallback
  const sys = gameSystems.find(gs => gs.id === systemId)
  return sys ? sys.code || fallback : fallback
}

export const systemCodeFor = (gameSystems = [], systemId, fallback = '') =>
  getSystemCodeById(gameSystems, systemId, fallback)

export const getSourceCode = (sources = [], sourceId) => {
  if (!sourceId) return ''
  const source = sources.find(src => src.id === sourceId)
  return source ? source.code || '' : ''
}

export const loadGameSystems = async () => {
  const resp = await API.get('/api/maintenance/game-systems')
  const systems = resp.data?.game_systems || []
  return systems.map(normalizeSystem)
}
