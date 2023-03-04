import React, {Component} from 'react';
import {Provider} from 'react-redux';
import {BrowserRouter} from 'react-router-dom';
import './App.css';
import SignIn from "../routes/signin/SignIn";
import store from "../redux/store";

class App extends Component {
    render() {
        // this.props.location.pathname
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <div>
                        <SignIn/>
                    </div>
                </BrowserRouter>
            </Provider>
        );
    }
}
export default App;
