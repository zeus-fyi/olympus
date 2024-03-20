import {Stack, Switch} from "@mui/material";

export const PageToggleView = (props: any) => {
    const { pageView, setPageView } = props;
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPageView(event.target.checked);
    };

    return (
        <div>
            <Stack direction={"row"} spacing={2} alignItems={"center"}>
            {pageView ? (
                <p>Advanced View</p>
                ) : (
                <p>Simplified View</p>
            )}
                <Switch
                    checked={pageView}
                    onChange={handleChange}
                    color="primary"
                    name="pageView"
                    inputProps={{ 'aria-label': 'toggle page view' }}
                />
            </Stack>
        </div>
    );
};