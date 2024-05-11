import React from "react";
import {
    BrowserRouter,
    Routes,
    Route,
} from "react-router-dom";
import UserForm from "./UserForm";
import "./App.css";
import Products from "./Products";
import ShoppingCart from "./Cart";

export default class App extends React.Component {
    state = {
        total: null,
        next: null,
        operation: null,
    };

    render() {
        return (
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<UserForm />} />
                    <Route path="/products" element={<Products />} />
                    <Route path="/cart" element={<ShoppingCart />} />
                </Routes>
            </BrowserRouter>
        );
    }
}