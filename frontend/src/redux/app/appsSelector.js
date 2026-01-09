const selectApps = (groups = {}, { query }) => {
  const keys = Object.keys(groups);
  let matchedGroups = {};

  keys.forEach(group => {
    const matchedGroup = groups[group].filter(app => {
      const nameMatch = app.name.toLowerCase().includes(query.toLowerCase());
      const groupMatch = app.group.toLowerCase().includes(query.toLowerCase());
      return nameMatch || groupMatch;
    });

    if (matchedGroup.length > 0) {
      matchedGroups[group] = matchedGroup;
    }
  });

  return matchedGroups;
};

export default selectApps;
