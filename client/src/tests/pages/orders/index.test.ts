import { fireEvent, render, waitFor, within } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import OrderList from '@/pages/Order/index.vue'
import { Mock } from 'vitest'
import { MOCK_ORDER_INFO } from '@/lib/mock/orderinfo'

global.fetch = vi.fn(() =>
  Promise.resolve({
    json: () => 
      Promise.resolve(MOCK_ORDER_INFO)
  })
) as Mock

const MOCK_PUSH = vi.fn()

const ResizeObserverMock = vi.fn(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}));

// Stub the global ResizeObserver
vi.stubGlobal('ResizeObserver', ResizeObserverMock);


vi.mock('vue-router',async () => {
  return {
    ...vi.importActual('vue-router'),
    useRoute: () => ({
      query: {
        limit: 5, 
        offset: 5,
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      },
      fullPath: '/orders?limit=5&offset=5&start_date=2024-01-01+12:00:00Z&end_date=2024-02-29+12:00:00Z'
    }),
    useRouter: () => ({
      push: MOCK_PUSH
    })
  }
})

const renderComponent = () => render(OrderList, {
  global: {
    mocks: {
      $router: {
        push: MOCK_PUSH
      }
    }
  }
})

const user = userEvent.setup()

describe('Orders page', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  test('Should render page correctly',async () => {
    const { getByTestId , getByRole } = renderComponent()

    await waitFor(() => {
      expect(global.fetch)
        .toHaveBeenLastCalledWith("http://localhost:8080/api/orders?limit=5&offset=5&start_date=2024-01-01+12:00:00Z&end_date=2024-02-29+12:00:00Z")
    })

    expect(within(getByTestId('date-range')).getByText("Jan 01, 2024 - Feb 29, 2024")).toBeDefined()
    
    // simulate sort
    expect(getByTestId('sort-desc')).toBeDefined()

    await fireEvent.click(getByTestId('sort-desc'))
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 5,
        offset: 5,
        sort_direction: "ASC",
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      }
    })

    // simulate search
    const searchInput = getByTestId('search-input')
    await user.type(searchInput, "product name{enter}")
    
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 5,
        offset: 5,
        search: "product name",
        sort_direction: "DESC",
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      }
    })

    // simulate change limit
    const limitDropdown = getByTestId('dropdown-page-limit')
    const buttonDropdown = within(limitDropdown).getByText('5/page')
    await fireEvent.click(buttonDropdown)

    const dropdownContent = getByRole('menu')
    await fireEvent.click(within(dropdownContent).getAllByRole('menuitem')[1])
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 10,
        offset: 0,
        sort_direction: "DESC",
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      }
    })

    // simulate go to next page
    const nextPageArrow = getByTestId('next-page-arrow')
    await fireEvent.click(nextPageArrow)
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 5,
        offset: 10,
        sort_direction: "DESC",
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      }
    })

    // simulate go to previous page
    const previous = getByTestId('previous-page-arrow')
    await fireEvent.click(previous)
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 5,
        offset: 0,
        sort_direction: "DESC",
        start_date: "2024-01-01 12:00:00Z",
        end_date: "2024-02-29 12:00:00Z"
      }
    })
  })
})