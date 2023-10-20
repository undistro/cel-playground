export function groupBy(data, by) {
  if (!by in data[0])
    throw new Error("You must provide a valid 'by' argument!");

  const reducedData = data.reduce((acc, cur) => {
    return {
      ...acc,
      [cur[by]]: [...(acc[cur[by]] ?? []), cur],
    };
  }, {});

  const groupedData = Object.entries(reducedData).map(([key, value]) => ({
    label: key,
    value,
  }));

  return groupedData;
}
