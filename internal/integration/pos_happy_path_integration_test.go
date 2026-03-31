package integration_test

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	domainorder "github.com/AXONcompany/POS/internal/domain/order"
	domainowner "github.com/AXONcompany/POS/internal/domain/owner"
	domainpayment "github.com/AXONcompany/POS/internal/domain/payment"
	domainproduct "github.com/AXONcompany/POS/internal/domain/product"
	domainsession "github.com/AXONcompany/POS/internal/domain/session"
	domaintable "github.com/AXONcompany/POS/internal/domain/table"
	domainuser "github.com/AXONcompany/POS/internal/domain/user"
	domainvenue "github.com/AXONcompany/POS/internal/domain/venue"
	usecaseauth "github.com/AXONcompany/POS/internal/usecase/auth"
	usecaseorder "github.com/AXONcompany/POS/internal/usecase/order"
	usecasepayment "github.com/AXONcompany/POS/internal/usecase/payment"
	usecaseproduct "github.com/AXONcompany/POS/internal/usecase/product"
	usecasetable "github.com/AXONcompany/POS/internal/usecase/table"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestPOSHappyPathIntegration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	state := newPOSState(map[int64]float64{
		9001: 1000, // azucar base para receta
	})

	userRepo := &authUserRepo{state: state}
	sessionRepo := &authSessionRepo{state: state}
	ownerRepo := &authOwnerRepo{state: state}
	venueRepo := &authVenueRepo{state: state}
	tableRepo := &posTableRepo{state: state}
	productRepo := &posProductRepo{state: state}
	categoryRepo := &posCategoryRepo{state: state}
	recipeRepo := &posRecipeRepo{state: state}
	orderRepo := &posOrderRepo{state: state}
	paymentRepo := &posPaymentRepo{state: state}

	authUC := usecaseauth.NewUsecase(userRepo, sessionRepo, "integration-secret", ownerRepo, venueRepo)
	tableUC := usecasetable.NewUsecase(tableRepo)
	productUC := usecaseproduct.NewUsecase(productRepo, categoryRepo, recipeRepo)
	orderUC := usecaseorder.NewUsecase(orderRepo, productRepo)
	paymentUC := usecasepayment.NewUsecase(paymentRepo)

	// 1) Auth + owner + venue
	ownerTokens, err := authUC.RegisterOwnerWithVenue(
		ctx,
		"Owner Uno",
		"owner@pos.dev",
		"strong-pass-123",
		"Sede Centro",
		"Calle 1",
		"3000000000",
		"it-test",
		"127.0.0.1",
	)
	require.NoError(t, err)
	require.NotEmpty(t, ownerTokens.AccessToken)
	require.Equal(t, 1, ownerTokens.User.RoleID)
	require.Equal(t, 1, ownerTokens.User.VenueID)

	parsed, err := jwt.Parse(ownerTokens.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("integration-secret"), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)

	waiter, waiterPassword, err := authUC.RegisterWaiter(ctx, "Mesero Uno", "mesero@pos.dev", ownerTokens.User.VenueID)
	require.NoError(t, err)
	require.Equal(t, 3, waiter.RoleID)
	require.NotEmpty(t, waiterPassword)

	waiterTokens, err := authUC.Login(ctx, waiter.Email, waiterPassword, "it-test", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, waiterTokens.AccessToken)

	// 2) Tables component
	tableEntity := &domaintable.Table{
		VenueID:  ownerTokens.User.VenueID,
		Number:   1,
		Capacity: 4,
		Status:   "libre",
	}
	err = tableUC.Create(ctx, tableEntity)
	require.NoError(t, err)
	require.NotZero(t, tableEntity.ID)

	occupied := "ocupada"
	err = tableUC.Update(ctx, tableEntity.ID, ownerTokens.User.VenueID, &domaintable.TableUpdates{Status: &occupied})
	require.NoError(t, err)

	loadedTable, err := tableUC.FindByID(ctx, tableEntity.ID, ownerTokens.User.VenueID)
	require.NoError(t, err)
	require.Equal(t, occupied, loadedTable.Status)

	// 3) Product/menu component
	createdCategory, err := productUC.CreateCategory(ctx, domainproduct.Category{
		VenueID: ownerTokens.User.VenueID,
		Name:    "Bebidas",
	})
	require.NoError(t, err)
	require.NotZero(t, createdCategory.ID)

	menuItem, err := productUC.CreateMenuItem(
		ctx,
		ownerTokens.User.VenueID,
		"Limonada",
		12.5,
		[]domainproduct.RecipeItem{
			{
				IngredientID:     9001,
				QuantityRequired: 50,
			},
		},
	)
	require.NoError(t, err)
	require.NotZero(t, menuItem.ID)

	// 4) Orders component (flujo normal)
	tableID := tableEntity.ID
	createdOrder, err := orderUC.CreateOrderWithoutItems(ctx, ownerTokens.User.VenueID, waiter.ID, &tableID)
	require.NoError(t, err)
	require.NotZero(t, createdOrder.ID)
	require.Equal(t, 1, createdOrder.StatusID)

	err = orderUC.AddProductToOrder(ctx, ownerTokens.User.VenueID, createdOrder.ID, []domainorder.OrderItem{
		{
			ProductID: menuItem.ID,
			Quantity:  2,
		},
	})
	require.NoError(t, err)

	updatedOrder, err := orderUC.GetOrderByID(ctx, ownerTokens.User.VenueID, createdOrder.ID)
	require.NoError(t, err)
	require.Len(t, updatedOrder.Items, 1)
	require.Equal(t, 2, updatedOrder.Items[0].Quantity)
	require.InDelta(t, 25.0, updatedOrder.TotalAmount, 0.0001)

	remainingStock := state.getIngredientStock(9001)
	require.InDelta(t, 900.0, remainingStock, 0.0001)

	err = orderUC.UpdateOrderStatus(ctx, ownerTokens.User.VenueID, updatedOrder.ID, 2) // enviada
	require.NoError(t, err)
	err = orderUC.UpdateOrderStatus(ctx, ownerTokens.User.VenueID, updatedOrder.ID, 4) // lista
	require.NoError(t, err)
	err = orderUC.CheckoutOrder(ctx, ownerTokens.User.VenueID, updatedOrder.ID) // pagada
	require.NoError(t, err)

	finalStatus, err := orderRepo.GetStatusByID(ctx, updatedOrder.ID, ownerTokens.User.VenueID)
	require.NoError(t, err)
	require.Equal(t, 5, finalStatus)

	// 5) Payment component
	payment, err := paymentUC.ProcessPayment(
		ctx,
		updatedOrder.ID,
		nil,
		"efectivo",
		updatedOrder.TotalAmount,
		2.5,
		"",
		ownerTokens.User.VenueID,
		waiter.ID,
	)
	require.NoError(t, err)
	require.NotZero(t, payment.ID)
	require.Equal(t, "aprobado", payment.Status)
	require.InDelta(t, 27.5, payment.Total, 0.0001)

	paymentsByOrder, err := paymentRepo.GetByOrderID(ctx, updatedOrder.ID)
	require.NoError(t, err)
	require.Len(t, paymentsByOrder, 1)

	invoice, err := paymentUC.GenerateInvoice(ctx, payment.ID)
	require.NoError(t, err)
	require.Equal(t, payment.ID, invoice["pago_id"])
	require.Equal(t, updatedOrder.ID, invoice["orden_id"])

	free := "libre"
	err = tableUC.Update(ctx, tableEntity.ID, ownerTokens.User.VenueID, &domaintable.TableUpdates{Status: &free})
	require.NoError(t, err)
	finalTable, err := tableUC.FindByID(ctx, tableEntity.ID, ownerTokens.User.VenueID)
	require.NoError(t, err)
	require.Equal(t, "libre", finalTable.Status)
}

type posState struct {
	mu sync.Mutex

	userSeq      int
	users        map[int]*domainuser.User
	userByEmail  map[string]int
	sessionSeq   int
	sessions     map[string]*domainsession.Session
	ownerSeq     int
	owners       map[int]*domainowner.Owner
	ownerByEmail map[string]int
	venueSeq     int
	venues       map[int]*domainvenue.Venue

	tableSeq int64
	tables   map[int64]*domaintable.Table

	categorySeq int64
	categories  map[int64]*domainproduct.Category
	productSeq  int64
	products    map[int64]*domainproduct.Product
	recipeSeq   int64
	recipes     map[int64][]domainproduct.RecipeItem

	ingredientStock map[int64]float64

	orderSeq     int64
	orderItemSeq int64
	orders       map[int64]*domainorder.Order

	paymentSeq int64
	payments   map[int64]*domainpayment.Payment
}

func newPOSState(initialStock map[int64]float64) *posState {
	stock := make(map[int64]float64, len(initialStock))
	for id, value := range initialStock {
		stock[id] = value
	}

	return &posState{
		users:           map[int]*domainuser.User{},
		userByEmail:     map[string]int{},
		sessions:        map[string]*domainsession.Session{},
		owners:          map[int]*domainowner.Owner{},
		ownerByEmail:    map[string]int{},
		venues:          map[int]*domainvenue.Venue{},
		tables:          map[int64]*domaintable.Table{},
		categories:      map[int64]*domainproduct.Category{},
		products:        map[int64]*domainproduct.Product{},
		recipes:         map[int64][]domainproduct.RecipeItem{},
		ingredientStock: stock,
		orders:          map[int64]*domainorder.Order{},
		payments:        map[int64]*domainpayment.Payment{},
	}
}

func (s *posState) now() time.Time {
	return time.Now().UTC()
}

func (s *posState) getIngredientStock(ingredientID int64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ingredientStock[ingredientID]
}

type authUserRepo struct{ state *posState }

func (r *authUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	id, ok := r.state.userByEmail[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return cloneUser(r.state.users[id]), nil
}

func (r *authUserRepo) GetByID(_ context.Context, id int) (*domainuser.User, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	u, ok := r.state.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return cloneUser(u), nil
}

func (r *authUserRepo) Create(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	if _, exists := r.state.userByEmail[u.Email]; exists {
		return nil, errors.New("email already exists")
	}

	r.state.userSeq++
	created := cloneUser(u)
	created.ID = r.state.userSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}
	created.UpdatedAt = created.CreatedAt

	r.state.users[created.ID] = cloneUser(created)
	r.state.userByEmail[created.Email] = created.ID

	return cloneUser(created), nil
}

func (r *authUserRepo) UpdateLastAccess(_ context.Context, id int) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	u, ok := r.state.users[id]
	if !ok {
		return errors.New("user not found")
	}

	now := r.state.now()
	u.LastAccess = &now
	u.UpdatedAt = now
	return nil
}

type authSessionRepo struct{ state *posState }

func (r *authSessionRepo) Create(_ context.Context, s *domainsession.Session) (*domainsession.Session, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.sessionSeq++
	created := cloneSession(s)
	created.ID = r.state.sessionSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}

	r.state.sessions[created.RefreshToken] = cloneSession(created)
	return cloneSession(created), nil
}

func (r *authSessionRepo) GetByToken(_ context.Context, refreshToken string) (*domainsession.Session, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	s, ok := r.state.sessions[refreshToken]
	if !ok {
		return nil, errors.New("session not found")
	}
	return cloneSession(s), nil
}

func (r *authSessionRepo) Revoke(_ context.Context, refreshToken string) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	s, ok := r.state.sessions[refreshToken]
	if !ok {
		return errors.New("session not found")
	}
	s.IsRevoked = true
	return nil
}

type authOwnerRepo struct{ state *posState }

func (r *authOwnerRepo) Create(_ context.Context, o *domainowner.Owner) (*domainowner.Owner, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	if _, exists := r.state.ownerByEmail[o.Email]; exists {
		return nil, errors.New("owner email already exists")
	}

	r.state.ownerSeq++
	created := cloneOwner(o)
	created.ID = r.state.ownerSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}
	created.UpdatedAt = created.CreatedAt

	r.state.owners[created.ID] = cloneOwner(created)
	r.state.ownerByEmail[created.Email] = created.ID
	return cloneOwner(created), nil
}

func (r *authOwnerRepo) GetByID(_ context.Context, id int) (*domainowner.Owner, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	o, ok := r.state.owners[id]
	if !ok {
		return nil, errors.New("owner not found")
	}
	return cloneOwner(o), nil
}

func (r *authOwnerRepo) GetByEmail(_ context.Context, email string) (*domainowner.Owner, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	id, ok := r.state.ownerByEmail[email]
	if !ok {
		return nil, errors.New("owner not found")
	}
	return cloneOwner(r.state.owners[id]), nil
}

type authVenueRepo struct{ state *posState }

func (r *authVenueRepo) Create(_ context.Context, v *domainvenue.Venue) (*domainvenue.Venue, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.venueSeq++
	created := cloneVenue(v)
	created.ID = r.state.venueSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}
	created.UpdatedAt = created.CreatedAt

	r.state.venues[created.ID] = cloneVenue(created)
	return cloneVenue(created), nil
}

type posTableRepo struct{ state *posState }

func (r *posTableRepo) Create(_ context.Context, t *domaintable.Table) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.tableSeq++
	now := r.state.now()
	created := cloneTable(t)
	created.ID = r.state.tableSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = now
	}
	if created.Status == "" {
		created.Status = "libre"
	}

	r.state.tables[created.ID] = cloneTable(created)
	*t = *cloneTable(created)
	return nil
}

func (r *posTableRepo) FindAll(_ context.Context, venueID int) ([]domaintable.Table, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	var out []domaintable.Table
	for _, t := range r.state.tables {
		if t.VenueID == venueID && t.DeletedAt == nil {
			out = append(out, *cloneTable(t))
		}
	}
	return out, nil
}

func (r *posTableRepo) FindByID(_ context.Context, id int64, venueID int) (*domaintable.Table, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	t, ok := r.state.tables[id]
	if !ok || t.VenueID != venueID || t.DeletedAt != nil {
		return nil, errors.New("table not found")
	}
	return cloneTable(t), nil
}

func (r *posTableRepo) Update(_ context.Context, id int64, venueID int, updates *domaintable.TableUpdates) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	t, ok := r.state.tables[id]
	if !ok || t.VenueID != venueID || t.DeletedAt != nil {
		return errors.New("table not found")
	}
	if updates.Number != nil {
		t.Number = *updates.Number
	}
	if updates.Capacity != nil {
		t.Capacity = *updates.Capacity
	}
	if updates.Status != nil {
		t.Status = *updates.Status
	}
	if updates.ArrivalTime != nil {
		arrival := *updates.ArrivalTime
		t.ArrivalTime = &arrival
	}
	now := r.state.now()
	t.UpdatedAt = &now
	return nil
}

func (r *posTableRepo) Delete(_ context.Context, id int64, venueID int) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	t, ok := r.state.tables[id]
	if !ok || t.VenueID != venueID || t.DeletedAt != nil {
		return errors.New("table not found")
	}
	now := r.state.now()
	t.DeletedAt = &now
	return nil
}

type posCategoryRepo struct{ state *posState }

func (r *posCategoryRepo) CreateCategory(_ context.Context, c domainproduct.Category) (*domainproduct.Category, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.categorySeq++
	created := c
	created.ID = r.state.categorySeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}
	r.state.categories[created.ID] = cloneCategory(&created)
	return cloneCategory(&created), nil
}

func (r *posCategoryRepo) GetByID(_ context.Context, id int64, venueID int) (*domainproduct.Category, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	c, ok := r.state.categories[id]
	if !ok || c.VenueID != venueID || c.DeletedAt != nil {
		return nil, errors.New("category not found")
	}
	return cloneCategory(c), nil
}

func (r *posCategoryRepo) GetAllCategories(_ context.Context, venueID int, _, _ int) ([]domainproduct.Category, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	var out []domainproduct.Category
	for _, c := range r.state.categories {
		if c.VenueID == venueID && c.DeletedAt == nil {
			out = append(out, *cloneCategory(c))
		}
	}
	return out, nil
}

func (r *posCategoryRepo) UpdateCategory(_ context.Context, c domainproduct.Category) (*domainproduct.Category, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	current, ok := r.state.categories[c.ID]
	if !ok || current.VenueID != c.VenueID || current.DeletedAt != nil {
		return nil, errors.New("category not found")
	}
	current.Name = c.Name
	now := r.state.now()
	current.UpdatedAt = &now
	return cloneCategory(current), nil
}

func (r *posCategoryRepo) DeleteCategory(_ context.Context, id int64, venueID int) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	current, ok := r.state.categories[id]
	if !ok || current.VenueID != venueID || current.DeletedAt != nil {
		return errors.New("category not found")
	}
	now := r.state.now()
	current.DeletedAt = &now
	return nil
}

type posProductRepo struct{ state *posState }

func (r *posProductRepo) CreateProduct(_ context.Context, p domainproduct.Product) (*domainproduct.Product, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	return r.createProductLocked(p), nil
}

func (r *posProductRepo) createProductLocked(p domainproduct.Product) *domainproduct.Product {
	r.state.productSeq++
	created := p
	created.ID = r.state.productSeq
	created.CreatedAt = r.state.now()
	r.state.products[created.ID] = cloneProduct(&created)
	return cloneProduct(&created)
}

func (r *posProductRepo) GetByID(_ context.Context, id int64, venueID int) (*domainproduct.Product, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	p, ok := r.state.products[id]
	if !ok || p.VenueID != venueID || p.DeletedAt != nil {
		return nil, errors.New("product not found")
	}
	return cloneProduct(p), nil
}

func (r *posProductRepo) GetAllProducts(_ context.Context, venueID int, _, _ int) ([]domainproduct.Product, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	var out []domainproduct.Product
	for _, p := range r.state.products {
		if p.VenueID == venueID && p.DeletedAt == nil {
			out = append(out, *cloneProduct(p))
		}
	}
	return out, nil
}

func (r *posProductRepo) UpdateProduct(_ context.Context, p domainproduct.Product) (*domainproduct.Product, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	current, ok := r.state.products[p.ID]
	if !ok || current.VenueID != p.VenueID || current.DeletedAt != nil {
		return nil, errors.New("product not found")
	}
	current.Name = p.Name
	current.SalesPrice = p.SalesPrice
	current.IsActive = p.IsActive
	now := r.state.now()
	current.UpdatedAt = &now
	return cloneProduct(current), nil
}

func (r *posProductRepo) DeleteProduct(_ context.Context, id int64, venueID int) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	current, ok := r.state.products[id]
	if !ok || current.VenueID != venueID || current.DeletedAt != nil {
		return errors.New("product not found")
	}
	now := r.state.now()
	current.DeletedAt = &now
	return nil
}

func (r *posProductRepo) CreateProductWithRecipe(_ context.Context, p domainproduct.Product, items []domainproduct.RecipeItem) (*domainproduct.Product, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	created := r.createProductLocked(p)
	for _, item := range items {
		r.state.recipeSeq++
		recipe := item
		recipe.ID = r.state.recipeSeq
		recipe.ProductID = created.ID
		r.state.recipes[created.ID] = append(r.state.recipes[created.ID], recipe)
		if _, exists := r.state.ingredientStock[recipe.IngredientID]; !exists {
			r.state.ingredientStock[recipe.IngredientID] = 0
		}
	}
	return cloneProduct(created), nil
}

func (r *posProductRepo) GetProductPrice(_ context.Context, productID int64, venueID int) (float64, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	p, ok := r.state.products[productID]
	if !ok || p.VenueID != venueID || p.DeletedAt != nil {
		return 0, errors.New("product not found")
	}
	return p.SalesPrice, nil
}

func (r *posProductRepo) GetRecipeLines(_ context.Context, productID int64) ([]domainorder.RecipeLine, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	items := r.state.recipes[productID]
	lines := make([]domainorder.RecipeLine, len(items))
	for i, item := range items {
		lines[i] = domainorder.RecipeLine{
			IngredientID:     item.IngredientID,
			QuantityRequired: item.QuantityRequired,
		}
	}
	return lines, nil
}

type posRecipeRepo struct{ state *posState }

func (r *posRecipeRepo) AddRecipeItem(_ context.Context, item domainproduct.RecipeItem) (*domainproduct.RecipeItem, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.recipeSeq++
	created := item
	created.ID = r.state.recipeSeq
	r.state.recipes[item.ProductID] = append(r.state.recipes[item.ProductID], created)
	return &created, nil
}

func (r *posRecipeRepo) GetByProductID(_ context.Context, productID int64) ([]domainproduct.RecipeItem, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	items := r.state.recipes[productID]
	out := make([]domainproduct.RecipeItem, len(items))
	copy(out, items)
	return out, nil
}

type posOrderRepo struct{ state *posState }

func (r *posOrderRepo) Create(_ context.Context, o *domainorder.Order) (*domainorder.Order, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.orderSeq++
	now := r.state.now()
	created := cloneOrder(o)
	created.ID = r.state.orderSeq
	if created.StatusID == 0 {
		created.StatusID = 1
	}
	created.CreatedAt = now
	created.UpdatedAt = now

	if len(created.Items) > 0 {
		created.TotalAmount = 0
		for i := range created.Items {
			r.state.orderItemSeq++
			created.Items[i].ID = r.state.orderItemSeq
			created.Items[i].OrderID = created.ID
			if created.Items[i].CreatedAt.IsZero() {
				created.Items[i].CreatedAt = now
			}
			created.Items[i].UpdatedAt = created.Items[i].CreatedAt
			created.TotalAmount += created.Items[i].UnitPrice * float64(created.Items[i].Quantity)
		}
	}

	r.state.orders[created.ID] = cloneOrder(created)
	return cloneOrder(created), nil
}

func (r *posOrderRepo) GetByID(_ context.Context, id int64, venueID int) (*domainorder.Order, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	o, ok := r.state.orders[id]
	if !ok || o.VenueID != venueID {
		return nil, errors.New("order not found")
	}
	return cloneOrder(o), nil
}

func (r *posOrderRepo) GetStatusByID(_ context.Context, id int64, venueID int) (int, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	o, ok := r.state.orders[id]
	if !ok || o.VenueID != venueID {
		return 0, errors.New("order not found")
	}
	return o.StatusID, nil
}

func (r *posOrderRepo) UpdateStatus(_ context.Context, id int64, venueID int, statusID int) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	o, ok := r.state.orders[id]
	if !ok || o.VenueID != venueID {
		return errors.New("order not found")
	}
	o.StatusID = statusID
	o.UpdatedAt = r.state.now()
	return nil
}

func (r *posOrderRepo) ListByTable(_ context.Context, tableID int64, venueID int) ([]domainorder.Order, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	var out []domainorder.Order
	for _, o := range r.state.orders {
		if o.VenueID == venueID && o.TableID != nil && *o.TableID == tableID {
			out = append(out, *cloneOrder(o))
		}
	}
	return out, nil
}

func (r *posOrderRepo) AddItemsWithInventory(_ context.Context, orderID int64, venueID int, items []domainorder.OrderItem, deductions []domainorder.StockDeduction) error {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	o, ok := r.state.orders[orderID]
	if !ok || o.VenueID != venueID {
		return errors.New("order not found")
	}

	for _, d := range deductions {
		available := r.state.ingredientStock[d.IngredientID]
		if available+1e-9 < d.Quantity {
			return domainorder.ErrInsufficientStock
		}
	}
	for _, d := range deductions {
		r.state.ingredientStock[d.IngredientID] -= d.Quantity
		if math.Abs(r.state.ingredientStock[d.IngredientID]) < 1e-9 {
			r.state.ingredientStock[d.IngredientID] = 0
		}
	}

	now := r.state.now()
	for _, it := range items {
		r.state.orderItemSeq++
		newItem := it
		newItem.ID = r.state.orderItemSeq
		newItem.OrderID = o.ID
		newItem.CreatedAt = now
		newItem.UpdatedAt = now

		o.Items = append(o.Items, newItem)
		o.TotalAmount += newItem.UnitPrice * float64(newItem.Quantity)
	}
	o.UpdatedAt = now
	return nil
}

type posPaymentRepo struct{ state *posState }

func (r *posPaymentRepo) Create(_ context.Context, p *domainpayment.Payment) (*domainpayment.Payment, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	r.state.paymentSeq++
	created := clonePayment(p)
	created.ID = r.state.paymentSeq
	if created.CreatedAt.IsZero() {
		created.CreatedAt = r.state.now()
	}
	r.state.payments[created.ID] = clonePayment(created)
	return clonePayment(created), nil
}

func (r *posPaymentRepo) GetByID(_ context.Context, id int64) (*domainpayment.Payment, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	p, ok := r.state.payments[id]
	if !ok {
		return nil, errors.New("payment not found")
	}
	return clonePayment(p), nil
}

func (r *posPaymentRepo) GetByOrderID(_ context.Context, orderID int64) ([]*domainpayment.Payment, error) {
	r.state.mu.Lock()
	defer r.state.mu.Unlock()

	var out []*domainpayment.Payment
	for _, p := range r.state.payments {
		if p.OrderID == orderID {
			out = append(out, clonePayment(p))
		}
	}
	return out, nil
}

func cloneUser(u *domainuser.User) *domainuser.User {
	if u == nil {
		return nil
	}
	cp := *u
	if u.Phone != nil {
		phone := *u.Phone
		cp.Phone = &phone
	}
	if u.LastAccess != nil {
		lastAccess := *u.LastAccess
		cp.LastAccess = &lastAccess
	}
	return &cp
}

func cloneSession(s *domainsession.Session) *domainsession.Session {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func cloneOwner(o *domainowner.Owner) *domainowner.Owner {
	if o == nil {
		return nil
	}
	cp := *o
	return &cp
}

func cloneVenue(v *domainvenue.Venue) *domainvenue.Venue {
	if v == nil {
		return nil
	}
	cp := *v
	return &cp
}

func cloneTable(t *domaintable.Table) *domaintable.Table {
	if t == nil {
		return nil
	}
	cp := *t
	if t.ArrivalTime != nil {
		arrival := *t.ArrivalTime
		cp.ArrivalTime = &arrival
	}
	if t.UpdatedAt != nil {
		updated := *t.UpdatedAt
		cp.UpdatedAt = &updated
	}
	if t.DeletedAt != nil {
		deleted := *t.DeletedAt
		cp.DeletedAt = &deleted
	}
	return &cp
}

func cloneCategory(c *domainproduct.Category) *domainproduct.Category {
	if c == nil {
		return nil
	}
	cp := *c
	if c.UpdatedAt != nil {
		updated := *c.UpdatedAt
		cp.UpdatedAt = &updated
	}
	if c.DeletedAt != nil {
		deleted := *c.DeletedAt
		cp.DeletedAt = &deleted
	}
	return &cp
}

func cloneProduct(p *domainproduct.Product) *domainproduct.Product {
	if p == nil {
		return nil
	}
	cp := *p
	if p.UpdatedAt != nil {
		updated := *p.UpdatedAt
		cp.UpdatedAt = &updated
	}
	if p.DeletedAt != nil {
		deleted := *p.DeletedAt
		cp.DeletedAt = &deleted
	}
	return &cp
}

func cloneOrder(o *domainorder.Order) *domainorder.Order {
	if o == nil {
		return nil
	}
	cp := *o
	if o.TableID != nil {
		tableID := *o.TableID
		cp.TableID = &tableID
	}
	if o.POSTerminalID != nil {
		posTerminalID := *o.POSTerminalID
		cp.POSTerminalID = &posTerminalID
	}
	if o.DeletedAt != nil {
		deleted := *o.DeletedAt
		cp.DeletedAt = &deleted
	}
	if len(o.Items) > 0 {
		cp.Items = make([]domainorder.OrderItem, len(o.Items))
		copy(cp.Items, o.Items)
	}
	return &cp
}

func clonePayment(p *domainpayment.Payment) *domainpayment.Payment {
	if p == nil {
		return nil
	}
	cp := *p
	if p.DivisionID != nil {
		divisionID := *p.DivisionID
		cp.DivisionID = &divisionID
	}
	if p.POSTerminalID != nil {
		posTerminalID := *p.POSTerminalID
		cp.POSTerminalID = &posTerminalID
	}
	return &cp
}
