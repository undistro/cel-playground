async function getExampleContentById(mode, exampleID) {
  const response = await fetch(
    `../../assets/examples/${mode.id}/${exampleID}.json`
  );
  const data = await response.json();
  return data;
}

export const ExampleService = {
  getExampleContentById,
};
