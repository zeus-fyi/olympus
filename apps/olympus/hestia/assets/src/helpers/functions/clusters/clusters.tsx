const R = require('ramda');

export const clustersFormat = R.curry((clusters: any) => {
    const convert = R.applySpec({
        cloudCtxNsID: R.prop('cloudCtxNsID'),
        cloudProvider: R.prop('cloudProvider'),
        context: R.prop('context'),
        region: R.prop('region'),
        namespace: R.prop('namespace'),
        createdAt: R.prop('createdAt'),
    });
    return convert(clusters);
});
