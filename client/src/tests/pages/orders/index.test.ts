import { fireEvent, render, waitFor } from '@testing-library/vue'
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


vi.mock('vue-router',async () => {
  return {
    ...vi.importActual('vue-router'),
    useRoute: () => ({
      query: {
        limit: 5, 
        offset: 0
      },
      fullPath: '/orders?limit=5&offset=0'
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

describe('Orders page', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  test('Should render page correctly',async () => {
    const { getByTestId } = renderComponent()

    await waitFor(() => {
      expect(global.fetch).toHaveBeenLastCalledWith("http://localhost:8080/api/orders?limit=5&offset=0")
    })

    expect(getByTestId('sort-desc')).toBeDefined()

    await fireEvent.click(getByTestId('sort-desc'))
    expect(MOCK_PUSH).toHaveBeenLastCalledWith({
      query: {
        limit: 5,
        offset: 0,
        sort_direction: "ASC"
      }
    })
  })
})