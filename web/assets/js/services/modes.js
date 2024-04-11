async function getModes() {
  const response = await fetch("../../assets/modes.json");
  const modes = await response.json();
  return modes;
}

export const ModesService = {
  getModes,
};
