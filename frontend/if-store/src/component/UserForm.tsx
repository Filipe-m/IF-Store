import React from 'react';
import { Navigate } from "react-router-dom";

export default class UserForm extends React.Component {
    state = {
        username: '',
        email: '',
        formSubmitted: false,
    };

    handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({username: event.target.value});
    };

    handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({email: event.target.value});
    };

    handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const {username, email} = this.state;

        try {
            const response = await fetch(`${process.env.REACT_APP_ACCOUNT_URL}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({username, email})
            });

            if (!response.ok) {
                new Error('error creating user');
            }

            const userData = await response.json();
            localStorage.setItem('userData', JSON.stringify(userData));

            this.setState({username: '', email: '', formSubmitted: true});
        } catch (error) {
            console.error(error);
        }
    };

    render() {
        const { formSubmitted } = this.state;

        if (formSubmitted) {
            return <Navigate to="/products" />;
        }

        return (
            <div className="container">
                <div className="d-flex justify-content-center align-items-center" style={{minHeight: '100vh'}}>
                    <form onSubmit={this.handleSubmit} style={{maxWidth: '400px', width: '100%'}}>
                        <div className="form-group">
                            <label htmlFor="username">Username:</label>
                            <input
                                type="text"
                                className="form-control"
                                id="username"
                                value={this.state.username}
                                onChange={this.handleUsernameChange}
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="email">Email:</label>
                            <input
                                type="email"
                                className="form-control"
                                id="email"
                                value={this.state.email}
                                onChange={this.handleEmailChange}
                            />
                        </div>
                        <div className="form-group" style={{marginTop: '1rem'}}>
                            <button type="submit" className="btn btn-primary btn-block">
                                Enviar
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        );
    }
}
