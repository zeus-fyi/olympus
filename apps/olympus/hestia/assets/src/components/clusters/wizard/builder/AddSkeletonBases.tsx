export function AddSkeletonBases(props: any) {
    // const {componentBase, componentBaseName} = props;
    //
    // const dispatch = useDispatch();
    //
    // let selectedComponentBaseSkeletonBasesKeys: string[] = []
    // if (componentBase !== null) {
    //     selectedComponentBaseSkeletonBasesKeys = Object.keys(componentBase);
    // }
    // const [inputField, setInputField] = useState('');
    //
    // const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    //     setInputField(event.target.value);
    // };
    //
    // const handleAddField = () => {
    //     if (inputField) {
    //         let sb = {  dockerImages: {},  };
    //         dispatch(addSkeletonBase({ componentBaseName: componentBaseName, skeletonBaseName: inputField, skeletonBase: sb }));
    //         setInputField('');
    //     }
    // };

    const handleRemoveField = (key: string) => {
        //dispatch(removeComponentBase(key));
    };

    return (
        <div>
            {/*{selectedComponentBaseSkeletonBasesKeys.map((key, index) => (*/}
            {/*    <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>*/}
            {/*        <TextField*/}
            {/*            fullWidth*/}
            {/*            id={`inputField-${index}`}*/}
            {/*            label={`Skeleton Base Name`}*/}
            {/*            variant="outlined"*/}
            {/*            value={key}*/}
            {/*            InputProps={{ readOnly: true }}*/}
            {/*            sx={{ flex: 1, mr: 2 }}*/}
            {/*        />*/}
            {/*        <Button variant="contained" sx={{ width: '100px' }} onClick={() => handleRemoveField(key)}>*/}
            {/*            Remove*/}
            {/*        </Button>*/}
            {/*    </Box>))*/}
            {/*}*/}
            {/*<Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>*/}
            {/*    <TextField*/}
            {/*        fullWidth*/}
            {/*        id="inputField-new"*/}
            {/*        label="New Skeleton Base Name"*/}
            {/*        variant="outlined"*/}
            {/*        value={inputField}*/}
            {/*        onChange={handleChange}*/}
            {/*        sx={{ flex: 1, mr: 2 }}*/}
            {/*    />*/}
            {/*    <Button variant="contained" sx={{ width: '100px' }} onClick={handleAddField}>*/}
            {/*        Add*/}
            {/*    </Button>*/}
            {/*</Box>*/}
        </div>
    )
}
