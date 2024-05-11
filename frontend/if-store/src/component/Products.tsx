import React from 'react';
import { Navigate, Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faShoppingCart } from '@fortawesome/free-solid-svg-icons';

interface Product {
    id: string;
    name: string;
    description: string;
    price: string;
    created_at: string;
}

interface ProductsState {
    products: Product[];
    loading: boolean;
    currentPage: number;
    showConfirmation: boolean;
}

export default class Products extends React.Component<{}, ProductsState> {
    constructor(props: {}) {
        super(props);
        this.state = {
            products: [],
            loading: true,
            currentPage: 1,
            showConfirmation: false,
        };
    }

    async componentDidMount() {
        this.fetchProducts(this.state.currentPage);
    }

    async fetchProducts(page: number) {
        try {
            const response = await fetch(`${process.env.REACT_APP_INVENTORY_URL}/product?limit=30&page=${page}`);
            if (!response.ok) {
                throw new Error('Erro ao carregar os produtos.');
            }
            const data = await response.json();
            this.setState({ products: data, loading: false});
        } catch (e) {
            console.error('Erro ao fazer a requisição:', e);
            this.setState({ loading: true});
        }
    }

    handlePreviousPage = () => {
        const { currentPage } = this.state;
        if (currentPage > 1) {
            this.setState({ currentPage: currentPage - 1 }, () => {
                this.componentDidMount();
            });
        }
    };

    handleNextPage = () => {
        const { currentPage } = this.state;
        this.setState({ currentPage: currentPage + 1 }, () => {
            this.componentDidMount();
        });
    };

    handleAddToCart = async (productId: string, quantity: number) => {
        try {
            const userData = localStorage.getItem('userData');
            if (!userData) {
                throw new Error('user not authenticated');
            }
            const userId = JSON.parse(userData).id;

            const response = await fetch(`${process.env.REACT_APP_ORDER_URL}/order-item`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'USER-ID': userId
                },
                body: JSON.stringify({
                    product_id: productId,
                    quantity: quantity
                })
            });

            if (!response.ok) {
                throw new Error('error adding product to cart');
            }

            this.setState({ showConfirmation: true });
            setTimeout(() => {
                this.setState({ showConfirmation: false });
            }, 2000);
        } catch (error) {
            console.error(error);
        }
    };

    render() {
        const { products, loading, currentPage, showConfirmation } = this.state;

        const userData = localStorage.getItem('userData');
        if (!userData) {
            return <Navigate to="/" />;
        }

        if (loading) {
            return <div>Carregando...</div>;
        }

        return (
            <div className="container">
                <div className="d-flex justify-content-between align-items-center mb-3">
                    <h2>Lista de Produtos</h2>
                    <Link to="/cart" className="btn btn-outline-primary">
                        <FontAwesomeIcon icon={faShoppingCart} />
                    </Link>
                </div>
                <div className="row">
                    {products && products.map((product) => (
                        <div className="col-md-4 mb-4" key={product.id}>
                            <div className="card">
                                <div className="card-body">
                                    <h5 className="card-title">{product.name}</h5>
                                    <p className="card-text">{product.description}</p>
                                    <p className="card-text">Preço: ${product.price}</p>
                                    <p className="card-text">Criado em: {new Date(product.created_at).toLocaleString()}</p>
                                    <div className="form-group">
                                        <input type="number" className="form-control" defaultValue="1" min="1" step="1" id={`quantity-${product.id}`} />
                                    </div>
                                    <button className="btn btn-primary" onClick={() => this.handleAddToCart(product.id, parseInt((document.getElementById(`quantity-${product.id}`) as HTMLInputElement).value))}>Adicionar ao Carrinho</button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
                <div className="d-flex justify-content-center">
                    <div className="pagination">
                        <button className="btn btn-primary" onClick={this.handlePreviousPage} disabled={currentPage === 1}>Anterior</button>
                        <span className="mx-2">Página {currentPage}</span>
                        <button className="btn btn-primary" onClick={this.handleNextPage}>Próxima</button>
                    </div>
                </div>
                <div className="alert alert-success fixed-bottom" role="alert" style={{ display: showConfirmation ? 'block' : 'none' }}>
                    Produto adicionado ao carrinho!
                </div>
            </div>
        );
    }
}
