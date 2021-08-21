public class Cell {

    public boolean isAlive() {
        return isAlive;
    }

    public void setAlive(boolean alive) {
        isAlive = alive;
    }

    public boolean Lives(int neighbors) {
        if (neighbors == 3) {
            return true;
        }
        if (neighbors == 2) {
            return isAlive;
        }
        return false;
    }

    private boolean isAlive;

    public Cell(boolean alive) {
        isAlive = alive;
    }

}
