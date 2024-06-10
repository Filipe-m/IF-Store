import React from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner } from '@fortawesome/free-solid-svg-icons';

interface CartItem {
    id: string;
    order_id: string;
    product_id: string;
    quantity: number;
    unit_price: string;
    name: string;
    created_at: string;
    updated_at: string;
}

interface Cart {
    id: string;
    user_id: string;
    status: string;
    total_amount: string;
    created_at: string;
    updated_at: string;
    items: CartItem[];
}

interface PaymentMethod {
    id: string;
    paymentType: string;
}

interface CartState {
    cart: Cart | null;
    loading: boolean;
    showConfirmation: boolean;
    loadingFinish: boolean;
    finishSuccess: boolean | null;
    showModal: boolean;
    address: string;
    neighborhood: string;
    city: string;
    zipCode: string;
    creditCardNumber: string;
    cardOwner: string;
    securityCode: string;
    expirationDate: string;
    activeTab: string;
    paymentMethods: PaymentMethod[];
    showAddCardForm: boolean;
    paymentMethodId: string;
}

export default class ShoppingCart extends React.Component<{}, CartState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            cart: null,
            loading: true,
            showConfirmation: false,
            loadingFinish: false,
            finishSuccess: null,
            showModal: false,
            address: '',
            neighborhood: '',
            city: '',
            zipCode: '',
            creditCardNumber: '',
            cardOwner: '',
            securityCode: '',
            expirationDate: '',
            activeTab: 'address',
            paymentMethods: [],
            showAddCardForm: false,
            paymentMethodId: ''
        };
    }

    async componentDidMount() {
        try {
            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_ORDER_URL}/order-actual`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'USER-ID': userId
                }
            });

            if (response.status === 404) {
                this.setState({ cart: null, loading: false });
                return;
            }

            if (!response.ok) {
                throw new Error('error fetching cart data');
            }

            const data = await response.json();
            this.setState({ cart: data, loading: false });
        } catch (error) {
            console.error(error);
            this.setState({ loading: true });
        }
    }

    handleRemoveItem = async (item: CartItem) => {
        try {
            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_ORDER_URL}/order-item`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'USER-ID': userId
                },
                body: JSON.stringify(item)
            });

            if (!response.ok) {
                throw new Error('error removing item from cart');
            }
            this.componentDidMount();
            this.setState({ showConfirmation: true });
            setTimeout(() => {
                this.setState({ showConfirmation: false });
            }, 2000);
        } catch (error) {
            console.error(error);
        }
    };

    handleFinishOrder = async () => {
        this.setState({ showModal: true });
        await this.fetchPaymentMethods();
    };

    fetchPaymentMethods = async () => {
        try {
            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_PAYMENT_URL}/paymentMethod/${userId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error('error fetching payment methods');
            }

            const paymentMethods = await response.json();
            console.log(paymentMethods);
            this.setState({ paymentMethods });
        } catch (error) {
            console.error(error);
        }
    };

    handleSubmitForm = async (event: React.FormEvent) => {
        event.preventDefault();

        try {
            this.setState({ loadingFinish: true, finishSuccess: null });

            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_ORDER_URL}/order/finish`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'USER-ID': userId
                },
                body: JSON.stringify({
                    order_id: this.state.cart?.id,
                    payment_method_id: this.state.paymentMethodId,
                    neighborhood: this.state.neighborhood,
                    city: this.state.city,
                    zipCode: this.state.zipCode,
                })
            });

            if (!response.ok) {
                throw new Error('error finishing order');
            }

            this.componentDidMount();
            this.setState({ finishSuccess: true, showModal: false });
        } catch (error) {
            console.error(error);
            this.setState({ finishSuccess: false });
        } finally {
            this.setState({ loadingFinish: false });
        }
    };

    handleAddCard = async (event: React.FormEvent) => {
        event.preventDefault();

        try {
            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_PAYMENT_URL}/paymentMethod/${userId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify([{
                    number: this.state.creditCardNumber,
                    card_holder: this.state.cardOwner,
                    expiration: this.state.expirationDate,
                    cvv: parseFloat(this.state.securityCode),
                }])
            });

            if (!response.ok) {
                throw new Error('error adding card');
            }

            await this.fetchPaymentMethods();
            this.setState({ showAddCardForm: false, creditCardNumber: '', cardOwner: '', securityCode: '', expirationDate: '' });
        } catch (error) {
            console.error(error);
        }
    };

    render() {
        const { cart, loading, showConfirmation, loadingFinish, finishSuccess, showModal, activeTab, paymentMethods, showAddCardForm } = this.state;

        if (loading) {
            return <div>Carregando carrinho...</div>;
        }

        if (!cart || cart.items.length === 0) {
            return <div>Carrinho vazio</div>;
        }

        return (
            <div className="container">
                <div className="alert alert-success fixed-top" role="alert" style={{ display: showConfirmation ? 'block' : 'none' }}>
                    Produto removido do carrinho!
                </div>
                <div className="d-flex justify-content-between align-items-center mb-3">
                    <h2>Carrinho de Compras</h2>
                    {cart.status === "pending" && !showModal && (
                        <button className="btn btn-primary" onClick={this.handleFinishOrder} disabled={loadingFinish}>
                            {loadingFinish ? <FontAwesomeIcon icon={faSpinner} spin /> : 'Finalizar Compra'}
                        </button>
                    )}
                    {showModal && (
                        <div className="modal" tabIndex={-1} role="dialog" style={{ display: showModal ? 'block' : 'none' }}>
                            <div className="modal-dialog modal-dialog-centered" role="document">
                                <div className="modal-content">
                                    <div className="modal-header">
                                        <h5 className="modal-title">Finalizar Compra</h5>
                                    </div>
                                    <div className="modal-body">
                                        <ul className="nav nav-tabs">
                                            <li className="nav-item">
                                                <button className={`nav-link ${activeTab === 'address' ? 'active' : ''}`} onClick={() => this.setState({ activeTab: 'address' })}>Endereço</button>
                                            </li>
                                            <li className="nav-item">
                                                <button className={`nav-link ${activeTab === 'payment' ? 'active' : ''}`} onClick={() => this.setState({ activeTab: 'payment' })}>Pagamento</button>
                                            </li>
                                        </ul>
                                        <form onSubmit={this.handleSubmitForm}>
                                            <div className="tab-content">
                                                <div className={`tab-pane fade ${activeTab === 'address' ? 'show active' : ''}`}>
                                                    <div className="form-group">
                                                        <label htmlFor="address">Endereço</label>
                                                        <input type="text" className="form-control" id="address" value={this.state.address} onChange={(e) => this.setState({ address: e.target.value })} required />
                                                    </div>
                                                    <div className="form-group">
                                                        <label htmlFor="neighborhood">Bairro</label>
                                                        <input type="text" className="form-control" id="neighborhood" value={this.state.neighborhood} onChange={(e) => this.setState({ neighborhood: e.target.value })} required />
                                                    </div>
                                                    <div className="form-group">
                                                        <label htmlFor="city">Cidade</label>
                                                        <input type="text" className="form-control" id="city" value={this.state.city} onChange={(e) => this.setState({ city: e.target.value })} required />
                                                    </div>
                                                    <div className="form-group">
                                                        <label htmlFor="zipCode">CEP</label>
                                                        <input type="text" className="form-control" id="zipCode" value={this.state.zipCode} onChange={(e) => this.setState({ zipCode: e.target.value })} required />
                                                    </div>
                                                </div>
                                                <div className={`tab-pane fade ${activeTab === 'payment' ? 'show active' : ''}`}>
                                                    <div className="form-group">
                                                        <label htmlFor="paymentMethods">Métodos de Pagamento Salvos</label>
                                                        <select className="form-control" id="paymentMethods" onChange={(e) => this.setState({ paymentMethodId: e.target.value })}>
                                                            <option value="">Selecione um método de pagamento</option>
                                                            {paymentMethods.map((method) => (
                                                                <option key={method.id} value={method.id}>
                                                                    {method.paymentType}
                                                                </option>
                                                            ))}
                                                        </select>
                                                    </div>
                                                    <button type="button" className="btn btn-secondary" onClick={() => this.setState({ showAddCardForm: !showAddCardForm })}>
                                                        {showAddCardForm ? 'Cancelar' : 'Adicionar Cartão'}
                                                    </button>
                                                    {showAddCardForm && (
                                                        <div>
                                                            <div className="form-group">
                                                                <label htmlFor="creditCardNumber">Número do Cartão de Crédito</label>
                                                                <input type="text" className="form-control" id="creditCardNumber" value={this.state.creditCardNumber} onChange={(e) => this.setState({ creditCardNumber: e.target.value })} required />
                                                            </div>
                                                            <div className="form-group">
                                                                <label htmlFor="cardOwner">Nome do Dono do Cartão</label>
                                                                <input type="text" className="form-control" id="cardOwner" value={this.state.cardOwner} onChange={(e) => this.setState({ cardOwner: e.target.value })} required />
                                                            </div>
                                                            <div className="form-group">
                                                                <label htmlFor="securityCode">Código de Segurança</label>
                                                                <input type="text" className="form-control" id="securityCode" value={this.state.securityCode} onChange={(e) => this.setState({ securityCode: e.target.value })} required />
                                                            </div>
                                                            <div className="form-group">
                                                                <label htmlFor="expirationDate">Data de Vencimento</label>
                                                                <input type="text" className="form-control" id="expirationDate" value={this.state.expirationDate} onChange={(e) => this.setState({ expirationDate: e.target.value })} required />
                                                            </div>
                                                            <button type="button" className="btn btn-primary" onClick={this.handleAddCard} disabled={loadingFinish}>
                                                                {loadingFinish ? <FontAwesomeIcon icon={faSpinner} spin /> : 'Salvar Cartão'}
                                                            </button>
                                                        </div>
                                                    )}
                                                </div>
                                            </div>
                                            <div className="modal-footer">
                                                <button type="button" className="btn btn-secondary" onClick={() => this.setState({ showModal: false })}>Fechar</button>
                                                <button type="submit" className="btn btn-primary" disabled={loadingFinish}>
                                                    {loadingFinish ? <FontAwesomeIcon icon={faSpinner} spin /> : 'Finalizar Compra'}
                                                </button>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
                {finishSuccess !== null && (
                    <div className={`alert ${finishSuccess ? 'alert-success' : 'alert-danger'}`} role="alert">
                        {finishSuccess ? 'Compra finalizada com sucesso!' : 'Erro ao finalizar a compra.'}
                    </div>
                )}
                <div>
                    <p>Status: {cart.status}</p>
                    <p>Total: ${cart.total_amount}</p>
                </div>
                <div className="row">
                    {cart.items.map((item) => (
                        <div className="col-md-4 mb-4" key={item.id}>
                            <div className="card">
                                <div className="card-body">
                                    <h5 className="card-title">{item.name}</h5>
                                    <p className="card-text">Quantidade: {item.quantity}</p>
                                    <p className="card-text">Preço Unitário: ${item.unit_price}</p>
                                    <p className="card-text">Preço Total: ${parseFloat(item.unit_price) * item.quantity}</p>
                                    {cart.status === "pending" && (
                                        <button className="btn btn-primary" onClick={() => this.handleRemoveItem(item)} disabled={loadingFinish}>
                                            {loadingFinish ? <FontAwesomeIcon icon={faSpinner} spin /> : 'Remover do Carrinho'}
                                        </button>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
                <Link to="/products" className="btn btn-primary">Voltar</Link>
            </div>
        );
    }
}
