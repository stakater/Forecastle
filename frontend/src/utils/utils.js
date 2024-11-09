export const getInitials = str => (str.length > 0 ? str[0] : "");

export const groupBy = key => array =>
  array.reduce((objectsByKeyValue, obj) => {
    const value = obj[key];
    objectsByKeyValue[value] = (objectsByKeyValue[value] || []).concat(obj);
    return objectsByKeyValue;
  }, {});

/* eslint-disable no-useless-escape */
export const isURL = str => {
  var regexp = /(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/;
  return regexp.test(str);
};

export const sortAlphabetically = (array, accessor = i => i) =>
  array.sort((a, b) => accessor(a).toLowerCase().localeCompare(accessor(b).toLowerCase(), 'en', {numeric: true}));
